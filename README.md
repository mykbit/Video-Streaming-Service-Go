# Video Streaming Service
[Features](#features) | [Building](#building) | [Running and Interacting](#running-and-interacting) | [Customization](#customization)

This project is an attempt to emulate the behavior of a video streaming service like Twitch. It is a distributed system that consists of producers a.k.a content creators, a server a.k.a broker and multiple clients. The server is responsible for forwarding the videos and metadata, and the clients are responsible for subscribing and unsubscribing to the producers they are interested in. The entities communicate with each other using UDP sockets, which allows for a fast and memory efficient communication. The system is designed to be fault tolerant, and it can handle a big number of clients and producers.

Docker is required to build and run the project. The project was tested on MacOS Sonoma 14.2.1. 

## Features
Video Streaming Service features a fully-implemented backend component, which can be observed in the following diagrams:

<p align="center">
  <img src="https://github.com/mykbit/Video-Streaming-Service-Go/assets/96201443/f143f76c-0a4a-433a-9d6a-c19b0e961436" alt="Docker Window Snapshot">
  <img src="https://github.com/mykbit/Video-Streaming-Service-Go/assets/96201443/40cb040b-901a-4671-ab47-ada42146ed32" alt="Wireshark Window Snapshot">
</p>

Other features include:
- Fault tolerance.
- Scalability.
- Fast and memory efficient communication.
- Subscription/Unsubscription system.
- Concurrent data processing and transferring.

## Building
Requirements:
 - Docker
 - Go 1.21 or higher

1. Clone the repository: `git clone https://github.com/mykbit/Video-Streaming-Service-Go.git`
2. Navigate to the project directory: `cd path/to/Video-Streaming-Service-Go`
3. Build the project: `docker compose up --build`

## Running and Interacting
Once you have successfully built the project, you can start interacting with the entities by taking control over clients. The following are the commands you should use in order to obtain the control:
1. Create a new terminal window.
2. Obtain manual control over a particular client: `docker attach <client_container_name>`

Once you have taken over a client, you can use the following commands to interact with the broker:
- `subscribe <producer_id>` - subscribe to a producer with a particular ID
- `unsubscribe <producer_id>` - unsubscribe from a producer with a particular ID

All the other commands will be ignored by the input parser.

Client container names are manually defined in the `docker-compose.yml` file. Producer ID's can be found in the `.env` file. 

## Customization
The project is highly flexible and can be easily customized. The following are the options you can change:
- The number of clients and producers in `docker-compose.yml`.
- Producer ID's in `.env`.
- Network and entities' addresses and ports in `.env`.
- Add new sample data to stream in `producer/` folder and specify the path to the data in `.env`.
- The delay duration which affects the start-up time of the producer in `.env`.
- The framerate of the video, defined as `rate`, in `.env`.
