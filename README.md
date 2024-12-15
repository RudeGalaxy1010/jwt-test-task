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
</p>2. ![image](https://github.com/user-attachments/assets/f9fd13f6-33e6-4fd2-aaa7-0e8d565d3c00)
</p>3. ![image](https://github.com/user-attachments/assets/d5e90be8-e84f-487b-abd2-24119d7f42f1)
</ul>
