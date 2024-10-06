if [ "$(docker network ls --filter name=e-commerce -q)" ]; then
    echo "E-commerce network exists."
else
    docker network create e-commerce
fi