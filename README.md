# AI DAG

This project presents a simple yet powerful implementation akin to a "poor man's" langgraph. As expected from a Directed Acyclic Graph (DAG), cyclic dependencies are not permitted. The core philosophy emphasizes leveraging straightforward tools that most users are already familiar with. Our approach involves managing a rudimentary DAG and traversing it in reverse, starting from the nodes with the fewest dependencies and progressing to the "root" of the DAG. Each step allows for the concurrent execution of independent computations, effectively making this a concurrent DAG.

In a demonstrative application, we utilize Google's Nearby Places and Weather APIs to get the weather and discover restaurants to consider visiting over the week. The decision on the optimal time and restaurants is then left to AI analysis. Despite its simplicity, the framework (defined in graph.yaml nodes) is designed to be adaptable, and I encourage collaboration on this project to expand the range of agents.

## Prerequisites

Before you begin, ensure you have met the following requirements:
- You have installed Go (version 1.xx or above) [Go Installation Guide](https://golang.org/doc/install)
- You are using a Linux or Mac OS machine. Windows is not tested but should work with appropriate environment variable setting methods.

## Setting Up Environment Variables

To use this project, you need to set up the following environment variables. Open your terminal and run:

```shell
export OPENAI_API_KEY='your_openai_apikey_here'
export OPEN_WEATHER_API_KEY='your_openweather_apikey_here' # currently only v3.0 supported, which need you to subscribe.
export GOOGLE_API_KEY='your_google_apikey_here'
```

Replace `your_openai_apikey_here`, `your_openweather_apikey_here`, and `your_google_apikey_here` with your actual API keys for OpenAI, OpenWeatherMap, and Google Cloud Services, respectively.

## Configuring the Graph.yaml

Before running the application, ensure you edit the `graph.yaml` to have the correct latitude and longitude for your area, or any other parameter you wish to tailor:

1. Open `graph.yaml` with your preferred text editor.
2. Change the `lat` and `lon` values to match your desired coordinates.
3. Feel free to adjust other settings in the file that you find relevant to your needs. This can include changing prompts or other configurations specific to this application.

Here's an example snippet you might modify:

```yaml
coordinates:
  lat: 12.3456
  lon: -65.4321
prompt: "Enter your customized prompt here if needed"
```

## Building the Project

To compile the project, navigate to the project directory in your terminal and run:

```shell
go build
```

This will generate an executable named `ai-dag` (or `ai-dag.exe` on Windows).

## Running the Application

After building the project, you can run it by executing:

```shell
./ai-dag
```

This will start the application using the configurations you've set. Make sure all previously mentioned setup steps have been correctly followed.

## Contributions

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

## License

[MIT](https://choosealicense.com/licenses/mit/)

## Like My Work and Want to Hire Me?
If you appreciate my work and are interested in hiring me, please feel free to reach out by clicking <a href="mailto:vinod@smartify.software">vinod@smartify.software</a> or contact me on LinkedIn at <a href="https://www.linkedin.com/in/vinod-halaharvi-289a1a13/">https://www.linkedin.com/in/vinod-halaharvi-289a1a13/</a>.
