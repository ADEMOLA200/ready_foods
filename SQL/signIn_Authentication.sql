{
    "email": "user@gmail.com",
    "password": "Password123"
}



curl -X GET \
  http://localhost:3000/login \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "user@gmail.com",
    "password": "Password123"
}'