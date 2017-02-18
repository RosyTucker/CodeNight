set -e

app="codenight-local"

echo ---- Removing existing codenight local containers ----

if docker ps | awk -v app="app" 'NR>1{  ($(NF) == app )  }'; then
  docker stop "$app" && docker rm -f "$app"
fi

godep save

echo ---- Building Image ----
docker build -t rosytucker/codenight .

echo ---- starting Container ----

source env.sh

echo Running on Port $PORT within container and port 4000 on host

winpty docker run -t --rm -i -e JWT_EXPIRY_HOURS=$JWT_EXPIRY_HOURS \
                        -e POST_LOGIN_REDIRECT=$POST_LOGIN_REDIRECT \
                        -e MONGODB_URI=$MONGODB_URI \
                        -e MASTER_USER=$MASTER_USER \
                        -e JWT_PUBLIC_KEY_BYTES="$JWT_PUBLIC_KEY_BYTES" \
                        -e JWT_PRIVATE_KEY_BYTES="$JWT_PRIVATE_KEY_BYTES" \
                        -e GITHUB_STATE_STRING=$GITHUB_STATE_STRING \
                        -e GITHUB_CALLBACK_URL=$GITHUB_CALLBACK_URL \
                        -e GITHUB_SECRET=$GITHUB_SECRET \
                         -e GITHUB_KEY=$GITHUB_KEY \
                        -e PORT=$PORT \
                        --name $app \
                        -p 127.0.0.1:4000:$PORT \
                        rosytucker/codenight
