package main

import "html/template"

var html = template.Must(template.New("chat_room").Parse(`
<html> 
<head> 
    <title>{{.roomid}}</title>
<style type="text/css">
html, body, div, span, applet, object, iframe,
h1, h2, h3, h4, h5, h6, p, blockquote, pre,
a, abbr, acronym, address, big, cite, code,
del, dfn, em, img, ins, kbd, q, s, samp,
small, strike, strong, sub, sup, tt, var,
b, u, i, center,
dl, dt, dd, ol, ul, li,
fieldset, form, label, legend,
table, caption, tbody, tfoot, thead, tr, th, td,
article, aside, canvas, details, embed, 
figure, figcaption, footer, header, hgroup, 
menu, nav, output, ruby, section, summary,
time, mark, audio, video {
	margin: 0;
	padding: 0;
	border: 0;
	font-size: 100%;
	font: inherit;
	vertical-align: baseline;
}
/* HTML5 display-role reset for older browsers */
article, aside, details, figcaption, figure, 
footer, header, hgroup, menu, nav, section {
	display: block;
}
body {
	line-height: 1;
}
ol, ul {
	list-style: none;
}
blockquote, q {
	quotes: none;
}
blockquote:before, blockquote:after,
q:before, q:after {
	content: '';
	content: none;
}
table {
	border-collapse: collapse;
	border-spacing: 0;
}
</style>
    <script src="//cdn.bootcss.com/jquery/1.7/jquery.min.js"></script>
    <script> 
        $('#message_form').focus();
        $(document).ready(function() { 
            // bind 'myForm' and provide a simple callback function 
            $('#myForm').on('submit',function(e) {
                e.preventDefault();
                $.post($(this).attr('action'),$(this).serializeArray(),function(){
                    $('#message_form').val('');
                    $('#message_form').focus();
                });
            });

            if (!!window.EventSource) {
                var source = new EventSource('/stream/{{.roomid}}');
                source.addEventListener('message', function(e) {
                    $('#messages').append(e.data + "</br>");
                    $('html, body').animate({scrollTop:$(document).height()}, 'slow');

                }, false);
            } else {
                alert("NOT SUPPORTED");
            }
        });
    </script> 
    </head>
    <body>
    <h1>Welcome to {{.roomid}} room</h1>
    <div id="messages"></div>
    <form id="myForm" action="/room/{{.roomid}}" method="post"> 
    User: <input id="user_form" name="user" value="{{.userid}}"></input> 
    Message: <input id="message_form" name="message"></input> 
    <input type="submit" value="Submit" /> 
    </form>
</body>
</html>
`))
