package templates

var Auth0UserHTML = `<html>
    <head>
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <script src="http://code.jquery.com/jquery-3.1.0.min.js" type="text/javascript"></script>

        <!-- font awesome from BootstrapCDN -->
        <link href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet">
        <link href="//maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css" rel="stylesheet">
        
        <script src="/public/js.cookie.js"></script>
        <script src="/public/user.js"> </script>
        <link href="/public/app.css" rel="stylesheet">
    </head>
    <body class="home">
        <div class="container">
            <div class="login-page clearfix">
              <div class="logged-in-box auth0-box logged-in">
                <h1 id="logo"><img src="/public/auth0_logo_final_blue_RGB.png" /></h1>
                <img class="avatar" src="{{.picture}}"/>
                <h2>Welcome {{.nickname}}</h2>
                <a id="qsLogoutBtn" class="btn btn-primary btn-lg btn-logout btn-block" href="/logout">Logout</a>
              </div>
            </div>
        </div>
    </body>
</html>`

var GoogleLoginButton = `<html>
<body>
	<a href="/login">Google Log In</a>
</body>
</html>`

var EmailTmpl = `To: {{namedAddresses .Mail.ToNames .Mail.To}}{{if .Mail.Cc}}
Cc: {{namedAddresses .Mail.CcNames .Mail.Cc}}{{end}}{{if .Mail.Bcc}}
Bcc: {{namedAddresses .Mail.BccNames .Mail.Bcc}}{{end}}
From: {{namedAddress .Mail.FromName .Mail.From}}
Subject: {{.Mail.Subject}}{{if .Mail.ReplyTo}}
Reply-To: {{namedAddress .Mail.ReplyToName .Mail.ReplyTo}}{{end}}
MIME-Version: 1.0
Content-Type: multipart/alternative; boundary="==============={{.Boundary}}=="
Content-Transfer-Encoding: 7bit
{{if .Mail.TextBody -}}
--==============={{.Boundary}}==
Content-Type: text/plain; charset=UTF-8
Content-Transfer-Encoding: 7bit
{{.Mail.TextBody}}
{{end -}}
{{if .Mail.HTMLBody -}}
--==============={{.Boundary}}==
Content-Type: text/html; charset=UTF-8
Content-Transfer-Encoding: 7bit
{{.Mail.HTMLBody}}
{{end -}}
--==============={{.Boundary}}==--
`
