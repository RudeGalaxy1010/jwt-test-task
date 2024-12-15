# jwt-test-task
Simple rest-api app to provide pair of tokens (JWT Access token + Refresh token)

<h1>Endpoints list:</h1>
<ul>
<li></p><h2>1. HTTP: /ping - route just to ping app.</h2>
</p>Request: None
</p>Response: "hello" - text
</li>
<li></p><h2>2. HTTP: /jwt - route to get pair of tokens.</h2>
</p>Rerquest: "guid" - user id (GUID)
</p>Response: {"access": "access_jwt_token", "refresh": "refresh_token"} - pair of tokens
</p>Note: Jwt token contains ip-adress of machine making request and refresh_token id
</li>
<li></p><h2>3. HTTP: /jwt-refresh - route to refresh and get new pair of tokens.</h2>
</p>Rerquest: {"access": "access_jwt_token", "refresh": "refresh_token"} - pair of tokens to refresh
</p>Response: {"access": "access_jwt_token", "refresh": "refresh_token"} - new pair of tokens
</li>
</ul>

<h1>Results:</h1>
<ul>
</p>2. Access token: eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzQyNzQyMDEsImlhdCI6MTczNDI3MzMwMSwiaXAiOiIxMjcuMC4wLjEiLCJqdGkiOiJlZDZkMzc5MS02MmMyLTRjYjktODEwNi1kYTcyMjg1NjcwZmYiLCJzdWIiOiJhc2NzIn0.Xp_Svc60X-xs_boqtcQohqnkG1fUXLUYMbKh1VkUbdOs-VwNNbc56PczJTR4CWZmed-m7yIZZHmNPOr2xXG5lw
</p> Refresh token: ZWQ2ZDM3OTEtNjJjMi00Y2I5LTgxMDYtZGE3MjI4NTY3MGZm
</p><image src ="https://github.com/user-attachments/assets/f9fd13f6-33e6-4fd2-aaa7-0e8d565d3c00"></image>
</p>3. Access token: eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzQyNzQyNzIsImlhdCI6MTczNDI3MzM3MiwiaXAiOiIxMjcuMC4wLjEiLCJqdGkiOiIyMWVmODNlMy1mM2UyLTRkYWMtYjc5Ni02OGEzZDk1MzU3MzciLCJzdWIiOiJhc2NzIn0.Au5Is_-9bh92rnNXisIB17lXIUof7ER7xBer6lifyhkB61JG4xei8LVPGnhDu36VYiMQJfuHJVAnahyzVOfugQ
</p> Refresh token: MjFlZjgzZTMtZjNlMi00ZGFjLWI3OTYtNjhhM2Q5NTM1NzM3
<image src ="https://github.com/user-attachments/assets/d5e90be8-e84f-487b-abd2-24119d7f42f1"></image>
</ul>
