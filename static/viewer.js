//For FPS
let frames = 0 ;
//Image storage stuff
let imageTypeSet = false ;
//Props
let width = 950 ;
let height = 600 ;
let currentSF ;
let paused = false ;
//Events
let actionQue = new Array() ;
let lastX ;
let lastY ;
let ctrlDown = false ;
//JS keys to X11 keysym
const keyTable = {
	'ArrowUp': 'Up',
	'ArrowDown': 'Down',
	'ArrowLeft': 'Left',
	'ArrowRight': 'Right',
	'PageUp': 'Prior',
	'PageDown': 'Next',
	'Enter': 'Return',
	'ContextMenu': 'Menu',
	'ScrollLock': 'Scroll_Lock',
	'Backspace': 'BackSpace',
	'!': 'exclam',
	' ': 'space',
	'"': 'quotedbl',
	'#': 'numbersign',
	'$': 'dollar',
	'%': 'percent',
	'&': 'ampersand',
	'\'': 'apostrophe',
	'(': 'parenleft',
	')': 'parenright',
	'*': 'asterisk',
	'+': 'plus',
	',': 'comma',
	'-': 'minus',
	'.': 'period',
	'/': 'slash',
	':': 'colon',
	';': 'semicolon',
	'<': 'less',
	'=': 'equal',
	'>': 'greater',
	'?': 'question',
	'@': 'at',
	'[': 'bracketleft',
	'\\': 'backslash',
	']': 'bracketright',
	'^': 'asciicircum',
	'_': 'underscore',
	'`': 'grave',
	'{': 'braceleft',
	'|': 'bar',
	'}': 'braceright',
	'~': 'asciitilde',
	'¡': 'exclamdown',
	'¢': 'cent ',
	'£': 'sterling',
	'¤': 'currency',
	'¥': 'yen',
	'¦': 'brokenbar',
	'§': 'section',
	'¨': 'diaeresis',
	'©': 'copyright',
	'ª': 'ordfeminine',
	'«': 'guillemotleft',
	'¬': 'notsign',
	'­': 'hyphen',
	'®': 'registered',
	'¯': 'macron',
	'°': 'degree',
	'±': 'plusminus',
	'²': 'twosuperior',
	'³': 'threesuperior',
	'´': 'acute',
	'µ': 'mu',
	'¶': 'paragraph',
	'·': 'periodcentered',
	'ç': 'cedilla',
	'¹': 'onesuperior',
	'º': 'masculine',
	'»': 'guillemotright',
	'¼': 'onequarter',
	'½': 'onehalf',
	'¾': 'threequarters',
	'¿': 'questiondown',
	'À': 'Agrave',
	'Á': 'Aacute',
	'Â': 'Acircumflex',
	'Ã': 'Atilde',
	'Ä': 'Adiaeresis',
	'Å': 'Aring',
	'Æ': 'AE',
	'Ç': 'Ccedilla',
	'È': 'Egrave',
	'É': 'Eacute',
	'Ê': 'Ecircumflex',
	'Ë': 'Ediaeresis',
	'Ì': 'Igrave',
	'Í': 'Iacute',
	'Î': 'Icircumflex',
	'Ï': 'Idiaeresis',
	'Ð': 'ETH',
	'Ñ': 'Ntilde',
	'Ò': 'Ograve',
	'Ó': 'Oacute',
	'Ô': 'Ocircumflex',
	'Õ': 'Otilde',
	'Ö': 'Odiaeresis',
	'×': 'multiply',
	'Ø': 'Ooblique',
	'Ù': 'Ugrave',
	'Ú': 'Uacute',
	'Û': 'Ucircumflex',
	'Ü': 'Udiaeresis',
	'Ý': 'Yacute',
	'Þ': 'THORN',
	'ß': 'ssharp',
	'à': 'agrave',
	'á': 'aacute',
	'â': 'acircumflex',
	'ã': 'atilde',
	'ä': 'adiaeresis',
	'å': 'aring',
	'æ': 'ae',
	'ç': 'ccedilla',
	'è': 'egrave',
	'é': 'eacute',
	'ê': 'ecircumflex',
	'ë': 'ediaeresis',
	'ì': 'igrave',
	'í': 'iacute',
	'î': 'icircumflex',
	'ï': 'idiaeresis',
	'ð': 'eth',
	'ñ': 'ntilde',
	'ò': 'ograve',
	'ó': 'oacute',
	'ô': 'ocircumflex',
	'õ': 'otilde',
	'ö': 'odiaeresis',
	'÷': 'division',
	'ø': 'oslash',
	'ù': 'ugrave',
	'ú': 'uacute',
	'û': 'ucircumflex',
	'ü': 'udiaeresis',
	'ý': 'yacute',
	'þ':'thorn',
	'ÿ': 'ydiaeresis'
} ;
//Takes and event, and a boolean (true of it is a keydown event, false if not) and returns the X11 keysym key for the given key
function getKey(e, down) {
	if (e.key === "Control") {
		//If it is the control key, set the ctrlDown variable
		ctrlDown = down ;
		//Return Control_L or Control_R, or Control if we don't know :(
		if (e.keyCode === "ControlLeft") {
			return "Control_L" ;
		}
		if (e.keyCode === "ControlRight") {
			return "Control_R" ;
		}
		return e.key ;
	} else if (!ctrlDown && e.ctrlKey) {
		//If the event says the control key is down, but we haven't registered it, press it!
		actionQue = actionQue.concat([4, String("Control_L").length], strToArray("Control_L")) ;
		ctrlDown = true ;
	} else if (ctrlDown && !e.ctrlKey) {
		//If the event says the control key is up, but we have registered it, unpress it!
		actionQue = actionQue.concat([5, String("Control_L").length], strToArray("Control_L")) ;
		ctrlDown = false ;
	}
	//If there is an entry in the ley table, return the entry
	if (keyTable[e.key]) {
		return keyTable[e.key] ;
	}
	//If it is a shift key, return Shift_L or Shift_R
	if (e.key === "Shift") {
		if (e.keyCode === "ShiftLeft") {
			return "Shift_L" ;
		}
		if (e.keyCode === "ShiftRight") {
			return "Shift_R" ;
		}
	}
	return e.key ;
}
//Updates the view with the latest frame and sends pending evnets
function updateFrame(view) {
	if (paused) {
		setTimeout(()=>updateFrame(view), 20) ;
		return ;
	}
	//Move the mouse to where it currently is
	actionQue = actionQue.concat([
		1, 255&(lastX>>24), 255&(lastX>>16), 255&(lastX>>8), 255&lastX, 255&(lastY>>24), 255&(lastY>>16), 255&(lastY>>8), 255&lastY
	]) ;
	
	fetch("/getframeanddo", {
		method: "POST",
		//Conver the actionQue to a buffer for the body
		body: (new Uint8Array(actionQue).buffer)
	}).then(resp=>{
		//Only carry on if we are OK
		if (resp.ok) {
			//If the type hasn't been set, set it to the content-type header
			if (!imageTypeSet) {
				view.type = resp.headers.get("content-type") ;
				imageTypeSet = true ;
			}
			//Get the data as a blob
			return resp.blob() ;
		}
		alert(`Server errored (${resp.status}).`) ;
		throw new Error(`Server responded with a ${resp.status}.`) ;
	}).then(image=>{
		//Remove the last image URL to save memory
		URL.revokeObjectURL(view.src) ;
		//Create a new URL and update the image
		view.src = URL.createObjectURL(image) ;
		frames++ ;
		//Start again...
		requestAnimationFrame(()=>updateFrame(view)) ;
	}) ;
	//Clear the actionQue because it has been sent now
	actionQue = [] ;
}
//Converts string to array of numbers representing the character set
function strToArray(s) {
	out = [] ;
	for (let c of s) {
		out.push(c.charCodeAt(0)) ;
	}
	return out ;
}
//Return the scale facter required to scale the image to the size of the window
function getScaleFactor() {
	let wSF = window.innerWidth / width ;
	let hSF = window.innerHeight / height ;
	return Math.min(wSF, hSF) ;
}
//Generate the new scale factor and resize the image
function updateSF(view) {
	currentSF = getScaleFactor() ;
	view.style.height = `${Math.floor(height*currentSF)}px` ;
	view.style.width = `${Math.floor(width*currentSF)}px` ;
}
//Start everything up
function start() {
	let view = document.getElementById("view") ;
	updateSF(view) ;
	
	//Mouse events
	view.addEventListener("mousedown", e=>{
		e.preventDefault() ;
		let x = Math.floor(e.offsetX / currentSF) ;
		let y = Math.floor(e.offsetY / currentSF) ;
		console.log(x, y) ;
		actionQue = actionQue.concat([1, 255&(x>>24), 255&(x>>16), 255&(x>>8), 255&x, 255&(y>>24), 255&(y>>16), 255&(y>>8), 255&y, 2, e.button+1]) ;
	}) ;
	view.addEventListener("mouseup", e=>{
		e.preventDefault() ;
		let x = Math.floor(e.offsetX / currentSF) ;
		let y = Math.floor(e.offsetY / currentSF) ;
		actionQue = actionQue.concat([1, 255&(x>>24), 255&(x>>16), 255&(x>>8), 255&x, 255&(y>>24), 255&(y>>16), 255&(y>>8), 255&y, 3, e.button+1]) ;
	}) ;
	//Save the current mouse position for later
	view.addEventListener("mousemove", e=>{
		lastX = Math.floor(e.offsetX / currentSF) ;
		lastY = Math.floor(e.offsetY / currentSF) ;
	}, {
		passive: true
	}) ;
	//Dont do default click events
	view.addEventListener("click", e=>e.preventDefault()) ;
	//Prevent context menu
	view.addEventListener("contextmenu", e=>e.preventDefault()) ;
	
	//Keyboard events
	window.addEventListener("keydown", e=>{
		e.preventDefault() ;
		tuKey = getKey(e, true) ;
		actionQue = actionQue.concat([4, tuKey.length], strToArray(tuKey)) ;
	}) ;
	window.addEventListener("keyup", e=>{
		e.preventDefault() ;
		tuKey = getKey(e, false) ;
		actionQue = actionQue.concat([5, tuKey.length], strToArray(tuKey)) ;
	}) ;
	
	//In console FPS meter
	setInterval(()=>{
		console.log(String(frames/5), 'fps') ;
		frames = 0 ;
	}, 5000) ;
	
	//Start the update loop
	updateFrame(view) ;
}
//RUN!!!
start() ;