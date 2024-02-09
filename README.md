# Project Golan
Creation of a graph of the number of new hospitalizations per year, by age and region.

## Folder
- data
    - file.csv (with data)
- graph 
    - image.png (image of the graph)
## Run Locally

Clone the project

```bash
  git clone https://github.com/TonyoCallimoutou/EFREI_TP_GO_graph
```

Go to the project directory

```bash
  cd EFREI_TP_GO_graph
```

Run with date filter

```bash 
  go run .\main.go -years="2021" -region="03" -age="59"
```

Run without date filter

```bash 
  go run .\main.go -region="03" -age="59"
```

## Screenshots

![graph](https://github.com/TonyoCallimoutou/EFREI_TP_GO_graph/blob/main/graph/hospitalisations_graph.png)
