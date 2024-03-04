{
    "otp": "094042"
}


curl -X POST \
  http://localhost:3000/verify-otp \
  -H 'Content-Type: application/json' \
  -d '{
    "otp": "001705"
}'