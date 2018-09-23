SWARM_MANAGER_IP=$1
DEPLOY_USER=travis-user

echo "Deploying to DockerHub"
echo "$DOCKERHUB_PASSWORD" | docker login -u "$DOCKERHUB_NAME" --password-stdin
docker tag swagger-ui-geo-rest $DOCKERHUB_NAME/swagger-ui-geo-rest
docker push $DOCKERHUB_NAME/swagger-ui-geo-rest

echo "Deploying to $SERVER_IP"

echo "Setting up ssh..."
eval "$(ssh-agent -s)"
ssh-keyscan -H SWARM_MANAGER_IP >> ~/.ssh/known_hosts
chmod 600 travis_key
ssh-add travis_key

echo "Uploading..."
scp -r docker-compose.yml $DEPLOY_USER@$SWARM_MANAGER_IP:/tmp/

echo "Pushing stack to swarm..."
ssh $DEPLOY_USER@$SWARM_MANAGER_IP "docker stack deploy --compose-file /tmp/docker-compose.yml geo-rest"