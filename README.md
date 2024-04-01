<p align="center">
  <img src="/images/AssetViz.png" alt="AssetViz">
</p>

# AssetViz
AssetViz simplifies the visualization of subdomains from input files, presenting them as a mind map. Ideal for penetration testers and bug bounty hunters conducting reconnaissance, AssetViz provides intuitive insights into domain structures for informed decision-making.

## Installation

Use Go 

```bash
go install -v github.com/mrrootsec/assetviz@latest
```
or

```bash
git clone https://github.com/mrrootsec/AssetViz
cd AssetViz
go install
```
## Usage

You can run this program using either of the following commands:

```bash
$ assetviz --help
Usage of assetviz:
  -f string
    	Path to the file containing subdomain names
```

```bash

assetviz -f filename
```
or

```bash

cat file | assetviz
```
This will generate output to the .report folder with the filename assetviz_report_date_time.html

## Screenshots
![AssetViz_1](/images/AssetViz_2.png)
![AssetViz_2](/images/AssetViz_1.png)

## License
This project is licensed under the MIT License

## Issue 

Feel free to open an issue if you have any problem with the script.
