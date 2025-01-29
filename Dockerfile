# Use a specific Go version to ensure consistency
FROM golang:1.22.3
# Set the working directory inside the container
WORKDIR /forum
# Copy the rest of the application code
COPY . .

RUN go mod download
# Build the application
RUN go build 

# Expose the application port
EXPOSE 8080
# Command to run the application
CMD [ "./forum" ]