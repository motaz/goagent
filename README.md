<h2>Introduction</h2>
GoAgent is a web service used by SimpleTrunk Panel web application to read and exceute commands in remote Asterisk PBX
https://github.com/motaz/stpanel
Written by Code for computer software (www.code.sd)

GoAgent executes as background service and listens to port 9091 in super user mode

<h2>Service Methods</h2>

<h3>Command</h3>
Executes CLI command in Asterisk<br/>
input parameters:<br/>
  
    Command

Example:</br>
{"command":"sip show peers"}<br/></br>

Output (JSON): <br/>

		success   bool  
		errorcode int   
		result    string
		message   string
                
<h3>Shell</h3>
Executes Linux shell command in Asterisk server<br/>
input parameter:<br/>
command<br/>
Example:<br/>
{"command":"uptime"}<br/></br>

Result (JSON):<br/>
		
    success   bool  
		errorcode int   
		result    string
		message   string
                
<h3>CallAMI</h3>
Executes AMI command in Asterisk server<br/>
Input parameter:<br/>
		Username string<br/>
		Secret   string<br/>
		Command  string<br/><br/>
   
output result (JSON):

		success   bool  
		errorcode int   
		result    string
		message   string

<h3>AddNode</h3>
Adds SIP node or Dialplan to specific configuratio file e.g.:<br/>

    [103]
    type=peer
    username=test
    secret=0987
    host=dynamic

Input parameters:
    
    filename string
		nodename string
		content  string
	
output result (JSON):

		success   bool  
		errorcode int   
		result    string
		message   string
		
		
