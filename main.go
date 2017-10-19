package main

// #cgo LDFLAGS: -lxdo
// #include <xdo.h>
// #include <stdio.h>
// #include <stdlib.h>
// #include <string.h>
/*
	//Stores the xdo type
	xdo_t* disp ;
	//Connect to the given display
	void init(const char* d) {
		disp = xdo_new(d) ;
	}
	//Tell the X server that the given keys are down
	void keysDown(const char* keys) {
		xdo_send_keysequence_window_down(disp, CURRENTWINDOW, keys, 0) ;
	}
	//Tell the X server that the given keys are up
	void keysUp(const char* keys) {
		xdo_send_keysequence_window_up(disp, CURRENTWINDOW, keys, 0) ;
	}
	//Tell the X server to move the mouse to the given coords
	void mouseTo(int x, int y) {
		xdo_move_mouse(disp, x, y, 0) ;
	}
	//Tell the X server the given mouse button is down
	void mouseDown(int button) {
		xdo_mouse_down(disp, CURRENTWINDOW, button) ;
	}
	//Tell the X server the given mouse button is up
	void mouseUp(int button) {
		xdo_mouse_up(disp, CURRENTWINDOW, button) ;
	}
	//Stores the length of the window list, getLengthOfList returns it
	int listLength ;
	int getLengthOfList() {
		return listLength ;
	}
	//Returns pointer to array of all windows
	Window *listWindows() {
		Window *list = NULL;

		//fill in search struct
        xdo_search_t search ;
		memset(&search, 0, sizeof(xdo_search_t)) ;
        search.winname = "." ;
		search.searchmask = SEARCH_NAME ;
		search.only_visible = 1 ;
		search.require = SEARCH_ALL ;
        search.limit = 0 ;
		search.max_depth = -1 ;

		xdo_search_windows(disp, &search, &list, &listLength) ;

		return list ;
	}
	char* getWindowName(unsigned long wi) {
		Window w = (Window)wi ;
		unsigned char* name ;
		int name_len ;
		int name_type ;
		xdo_get_window_name(disp, w, &name, &name_len, &name_type) ;
		return (char*)name ;
	}
	int getWindowPID(unsigned long wi) {
		Window w = (Window)wi ;
		return xdo_get_pid_window(disp, w) ;
	}
	void minWindow(unsigned long wi) {
		Window w = (Window)wi ;
		xdo_minimize_window(disp, w) ;
	}
	void closeWindow(unsigned long wi) {
		Window w = (Window)wi ;
		xdo_close_window(disp, w) ;
	}
	void killWindow(unsigned long wi) {
		Window w = (Window)wi ;
		xdo_kill_window(disp, w) ;
	}
	unsigned int *getWindowSize(unsigned long wi) {
		Window w = (Window)wi ;
		unsigned int *d ;
		d = (unsigned int *)malloc(4) ;
		xdo_get_window_size(disp, w, &d[0], &d[1]) ;
		return d ;
	}
	void setWindowSize(unsigned long wi, int w, int h) {
		Window win = (Window)wi ;
		xdo_set_window_size(disp, win, w, h, 0) ;
	}
	int *getWindowLocation(unsigned long wi) {
		Window w = (Window)wi ;
		int *p ;
		p = (int *)malloc(4) ;
		Screen *s ;
		xdo_get_window_location(disp, w, &p[0], &p[1], &s) ;
		return p ;
	}
	void setWindowLocation(unsigned long wi, int x, int y) {
		Window w = (Window)wi ;
		xdo_move_window(disp, w, x, y) ;
	}
*/
import "C"
import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	display   string  = ":2056"
	httpPort  string  = ":8080"
	width     int     = 950
	height    int     = 600
	depth     int     = 16
	scaler    float64 = 1
	imageType string  = "jpg"
	quality   string  = "100"
)

var startPath string

type handlerType struct{}

func (f handlerType) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	//Default headers
	resp.Header().Set("Server", "xrc")
	resp.Header().Set("Content-Type", "text/plain")
	//If we need to serve a static file
	if strings.Index(req.URL.Path, "/static/") == 0 {
		loadPath := filepath.Join(startPath, "..", req.URL.Path)
		//Make sure we are in the static directory
		if strings.Index(loadPath, startPath) != 0 {
			resp.WriteHeader(403)
			resp.Write([]byte("Forbidden"))
		} else {
			//Set the mime type header
			t := mime.TypeByExtension(filepath.Ext(loadPath))
			resp.Header().Set("Content-Type", t)
			//Serve the file
			http.ServeFile(resp, req, loadPath)
		}
	} else {
		//Get the page handler
		toCall, ok := pages[req.URL.Path]
		if ok {
			//Call the handler if it's available
			toCall(resp, req)
		} else {
			//Send default page
			resp.Header().Set("Content-Type", "text/html")
			resp.WriteHeader(200)
			resp.Write([]byte(`<html><head><title>xrc</title></head><body><h1 style="padding: 3px; margin: 0px;">xrc</h1><a href="/static/viewer.html">Viewer</a><br><a href="/static/controler.html">Controler</a><hr><i>Jacob O'Toole. MIT Licence.</i></body></html>`))
		}
	}
}

func getFrame(resp http.ResponseWriter, _ *http.Request) {
	//Basic argu,ents
	args := []string{"-display", display, "-window", "root", "-delay", "0", "-quality", quality}
	//If we need to resize the image, add the resize argument
	if scaler != 1 {
		imageWidth := int(math.Ceil(float64(width) * scaler))
		imageHeight := int(math.Ceil(float64(height) * scaler))
		args = []string{"-resize", strconv.Itoa(imageWidth) + "x" + strconv.Itoa(imageHeight)}
	}
	//Make the command with the image type added
	importCommand := exec.Command("import", append(args, imageType+":-")...)
	//Pipe stdin and error
	importCommand.Stdin = os.Stdin
	importCommand.Stderr = os.Stderr
	//Get the stdout output
	image, err := importCommand.Output()
	if err != nil {
		//Send a 500 and log the error if there was one
		resp.WriteHeader(500)
		fmt.Fprintln(os.Stderr, err)
		resp.Write([]byte("Internal Server Error"))
	} else {
		resp.WriteHeader(200)
		//Set the content-type
		resp.Header().Set("Content-Type", "image/"+imageType)
		//Write the data
		resp.Write(image)
	}
}

//Stores basic data about a window
type winInfo struct {
	ID       uint32
	Title    string
	Size     [2]uint64
	Location [2]int64
}

func getWins(resp http.ResponseWriter, _ *http.Request) {
	//Get the windows
	wins := listWindows()

	//Create a JSON object from the window list
	var jsonWins []interface{}
	for _, win := range wins {
		jsonWins = append(jsonWins, winInfo{
			win,
			C.GoString(C.getWindowName(C.ulong(win))),
			getWindowSize(win),
			getWindowLocation(win),
		})
	}
	data, err := json.Marshal(interface{}(jsonWins))
	if err != nil {
		resp.WriteHeader(500)
		fmt.Fprintln(os.Stderr, err)
		resp.Write([]byte("Internal Server Error"))
	} else {
		resp.WriteHeader(200)
		resp.Header().Set("Content-Type", "text/json")
		resp.Write(data)
	}
}

//Reads the body from the request & returns it
func readBody(req *http.Request) []byte {
	var out []byte
	buff := make([]byte, 1024)
	var n int
	var err error
	for {
		n, err = req.Body.Read(buff)
		out = append(out, buff[:n]...)
		if err != nil {
			if err == io.EOF {
				return out
			}
			panic(err)
		}
	}
}

//Represent window ID and position/size
type opts struct {
	win  C.ulong
	opts []C.int
}

//Parse data in the form winid-x-y to opts struct
func parsePosOpts(data []byte) opts {
	//Split the data
	optsBP := strings.Split(string(data), "-")
	//Get win
	win, err := strconv.ParseUint(optsBP[0], 10, 32)
	if err != nil {
		panic(err)
	}
	//Get x/y // w/h
	a1, err := strconv.Atoi(optsBP[1])
	if err != nil {
		panic(err)
	}
	a2, err := strconv.Atoi(optsBP[2])
	if err != nil {
		panic(err)
	}
	return opts{
		C.ulong(win),
		[]C.int{C.int(a1), C.int(a2)},
	}
}

//Sends a basic HTTP 200 response with the body provided
func sendBasic200(resp http.ResponseWriter, body string) {
	resp.Header().Set("Content-Length", strconv.Itoa(len(body)))
	resp.WriteHeader(200)
	resp.Write([]byte(body))
}

//Resizes window
func resizeWindow(resp http.ResponseWriter, req *http.Request) {
	o := parsePosOpts(readBody(req))
	C.setWindowSize(o.win, o.opts[0], o.opts[1])
	sendBasic200(resp, "")
}

//Moves window
func moveWindow(resp http.ResponseWriter, req *http.Request) {
	o := parsePosOpts(readBody(req))
	C.setWindowLocation(o.win, o.opts[0], o.opts[1])
	sendBasic200(resp, "Moved")
}

var interactCommands = map[byte]func(){
	//Pause
	0: func() {
		pause := binary.BigEndian.Uint32(commands[i : i+4])
		i += 4
		time.Sleep(time.Millisecond * time.Duration(pause))
	},
	//Mouse move  x:4, y:4
	1: func() {
		x := binary.BigEndian.Uint32(commands[i : i+4])
		i += 4
		y := binary.BigEndian.Uint32(commands[i : i+4])
		i += 4
		C.mouseTo(C.int(int(x)), C.int(int(y)))
	},
	//Mouse Down  button:1
	2: func() {
		C.mouseDown(C.int(int(commands[i])))
		i++
	},
	//Mouse Up  button:1
	3: func() {
		C.mouseUp(C.int(int(commands[i])))
		i++
	},
	//Key Down  length:1 char:<length>
	4: func() {
		len := int(commands[i])
		i++
		str := string(commands[i:(i + len)])
		i += len
		C.keysDown(C.CString(str))
	},
	//Key Up  length:1 char:<length>
	5: func() {
		len := int(commands[i])
		i++
		str := string(commands[i:(i + len)])
		i += len
		C.keysUp(C.CString(str))
	},
}

//Stores commands
var commands []byte

//Is a command parser loop running?
var commandLoopRunning = false

//Stores current command index
var i int

//Starts the command parser
func command() {
	commandLoopRunning = true
	var tc func()
	var ok bool
	for i < len(commands) {
		tc, ok = interactCommands[commands[i]]
		i++
		if ok {
			tc()
		}
	}
	commandLoopRunning = false
}

//Adds the commands to the command list and runs the parser if it isn't running
func runCommand(resp http.ResponseWriter, req *http.Request) {
	commands = append(commands, readBody(req)...)
	if !commandLoopRunning {
		go command()
	}
	sendBasic200(resp, "Recieved")
}

//Same as runCommand being run before getFrame
func getFrameAndDo(resp http.ResponseWriter, req *http.Request) {
	commands = append(commands, readBody(req)...)
	if !commandLoopRunning {
		go command()
	}
	getFrame(resp, req)
}

var pages = map[string]func(http.ResponseWriter, *http.Request){
	"/getframe":      getFrame,
	"/getwins":       getWins,
	"/resize":        resizeWindow,
	"/move":          moveWindow,
	"/do":            runCommand,
	"/getframeanddo": getFrameAndDo,
}

var xvfb xvfbProcess

type xvfbProcess struct {
	display string
	width   int
	height  int
	depth   int
	process *exec.Cmd
	running bool
}

//Starts Xvfb and waits to make sure it is loaded properly
func (p *xvfbProcess) start() bool {
	if p.running {
		return false
	}
	p.process = exec.Command("Xvfb", "-screen", "0", strconv.Itoa(p.width)+"x"+strconv.Itoa(p.height)+"x"+strconv.Itoa(p.depth), p.display)
	p.process.Stdin = os.Stdin
	p.process.Stdout = os.Stdout
	p.process.Stderr = os.Stderr
	p.process.Start()
	p.running = true
	time.Sleep(5 * time.Second)
	return true
}

func (p *xvfbProcess) stop() {
	p.process.Process.Kill()
	p.running = false
}

func startXvfb(display string, width, height, depth int) xvfbProcess {
	this := xvfbProcess{display, width, height, depth, nil, false}
	this.start()
	return this
}

//Returns a slice of window IDs from all the windows on the X server
func listWindows() []uint32 {
	//Get a list of windows
	winList := C.listWindows()
	//This is the pointer to the first item
	listAddr := uintptr(unsafe.Pointer(winList))
	//This is the length of each item in the list
	itemLength := unsafe.Sizeof(*winList)
	//And this is how many items there are
	lengthOfList := int(C.getLengthOfList())
	var list []uint32
	//For each item, convert it to a uint32 and append it
	for i := 0; i < lengthOfList; i++ {
		//                  Converts a pointer to the window ID to a pointer of a uint32 and then redirects it
		list = append(list, *(*uint32)(unsafe.Pointer(listAddr + uintptr(i)*itemLength)))
	}
	return list
}

//Returns size of a window win in the form [w, h]
func getWindowSize(win uint32) [2]uint64 {
	winSize := C.getWindowSize(C.ulong(win))
	winSizeAddr := unsafe.Pointer(winSize)
	defer C.free(winSizeAddr)
	itemLength := unsafe.Sizeof(*winSize)
	if itemLength == 4 {
		return [2]uint64{
			uint64(*(*uint32)(winSizeAddr)),
			uint64(*(*uint32)(unsafe.Pointer(uintptr(winSizeAddr) + itemLength))),
		}
	}
	return [2]uint64{
		*(*uint64)(winSizeAddr),
		*(*uint64)(unsafe.Pointer(uintptr(winSizeAddr) + itemLength)),
	}
}

//Returns location of a window win in the form [x, y]
func getWindowLocation(win uint32) [2]int64 {
	winLoc := C.getWindowLocation(C.ulong(win))
	winLocAddr := unsafe.Pointer(winLoc)
	defer C.free(winLocAddr)
	itemLength := unsafe.Sizeof(*winLoc)
	if itemLength == 4 {
		return [2]int64{
			int64(*(*int32)(winLocAddr)),
			int64(*(*int32)(unsafe.Pointer(uintptr(winLocAddr) + itemLength))),
		}
	}
	return [2]int64{
		*(*int64)(winLocAddr),
		*(*int64)(unsafe.Pointer(uintptr(winLocAddr) + itemLength)),
	}
}

func main() {
	fmt.Println("Loading xrc")
	//Get the dir path of the executable
	var err error
	startPath, err = os.Executable()
	if err != nil {
		panic(err)
	}
	startPath = filepath.Join(filepath.Dir(startPath), "static")
	//Start Xvfb
	fmt.Println("Starting Xvfb...")
	xvfb = startXvfb(display, width, height, depth)
	fmt.Println("Display server open on", display)
	//Init libxdo
	C.init(C.CString(display))
	fmt.Println("Hopefully xdo is connected to display", display)
	//Start the http server
	fmt.Println("Starting HTTP server on", httpPort)
	go http.ListenAndServe(httpPort, handlerType{})
	//Run bash with command line arguments
	fmt.Println("Running", "bash", strings.Join(os.Args[1:], " ")+"...")
	bash := exec.Command("bash", os.Args[1:]...)
	//Add DISPLAY environment variable
	bash.Env = append(os.Environ(), "DISPLAY="+display)
	//Pipe stderr, in and out and run
	bash.Stdin = os.Stdin
	bash.Stdout = os.Stdout
	bash.Stderr = os.Stderr
	err = bash.Run()
	//Kill xvfb and exit when we're done
	xvfb.stop()
	if err == nil {
		os.Exit(0)
	} else {
		//Exit with the exit code of bash
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			} else {
				fmt.Fprintln(os.Stderr, "bash exit status unknown, exit status 128")
				os.Exit(128)
			}
		}
	}
}
