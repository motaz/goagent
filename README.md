<h2>Introduction</h2>
GoAgent is a web service used by SimpleTrunk Panel web application to read and exceute commands in remote Asterisk PBX
https://github.com/motaz/stpanel
Written by Code for computer software (www.code.sd)

GoAgent executes as background service and listens to port 9091 in super user mode

<h2>Service Methods</h2>

<h3>Command</h3>
Executes CLI command in Asterisk
input parameters:
Command

Example:
{"command":"sip show peers"}

Output (JSON):
		success   bool   
		errorcode int    
		result    string 
		message   string 
