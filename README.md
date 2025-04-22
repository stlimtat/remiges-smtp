<a id="readme-top"></a>
<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![Unlicense License][license-shield]][license-url]


<!-- PROJECT LOGO -->
# remiges-smtp
SMTP client with file scraping capabilities


<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#tutorial">Tutorial</a></li>
    <li><a href="#faq">FAQ</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>


<!-- DIRECTORY STRUCTURE -->
## Directory structure
```
/remiges-smtp
├── config                 # Location to be used to store
├── docker-compose.yml     # Docker Compose configuration for smtpclientd
├── Dockerfile.smtpclientd # Dockerfile that compiles the application
├── output                 # Location to be used to record output
├── README.md              # Default README.md
├── testdata               # Location to be used to hold data
└── doc                    # Document directory
  ├── FAQ.md               # Frequently Asked Questions
  ├── QUICKSTART.md        # QuickStart Guide
  ├── TUTORIAL.md          # Tutorial Document
  └── USAGE.md             # Document providing usage information
```

<!-- ABOUT THE PROJECT -->
## About The Project

Remiges SMTP is a powerful SMTP client that can:
- Read and process files from directories
- Format them as RFC-compliant emails
- Sign messages with DKIM
- Connect to mail servers and send emails
- Record results in CSV format

### Built With

The following software were used to build this project and application:
- [Bazel](https://bazel.build)
  Bazel is used to build and run the application.
- [Cobra](https://github.com/spf13/cobra)
- [Viper](https://github.com/spf13/viper)
- [Zerolog](https://github.com/rs/zerolog)
- [mox](https://github.com/mjl-/mox)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

See [QUICKSTART.md](./doc/QUICKSTART.md)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/stlimtat/remiges-smtp.git
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- USAGE EXAMPLES -->
## Usage

See [USAGE.md](./doc/USAGE.md)

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- TUTORIAL -->
## Tutorial

See [TUTORIAL.md](./doc/TUTORIAL.md)

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- FAQ -->
## FAQ

See [FAQ.md](./doc/FAQ`.md)

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- ROADMAP -->
## Roadmap

- [x] Write initial smtp client
- [x] Provide gen dkim tool that helps provide DNS TXT entry and config
- [ ] Update README to allow users to use

See the [open issues](https://github.com/stlimtat/remiges-smtp/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

We welcome contributions! Please see our [Contributing Guide](./doc/CONTRIBUTING.md) for details.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the MIT License. See [LICENSE.MIT](./LICENSE.MIT) for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Project Link: [https://github.com/stlimtat/remiges-smtp](https://github.com/stlimtat/remiges-smtp)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ACKNOWLEDGMENTS -->
## Acknowledgments


<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/stlimtat/remiges-smtp.svg?style=for-the-badge
[contributors-url]: https://github.com/stlimtat/remiges-smtp/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/stlimtat/remiges-smtp.svg?style=for-the-badge
[forks-url]: https://github.com/stlimtat/remiges-smtp/network/members
[stars-shield]: https://img.shields.io/github/stars/stlimtat/remiges-smtp.svg?style=for-the-badge
[stars-url]: https://github.com/stlimtat/remiges-smtp/stargazers
[issues-shield]: https://img.shields.io/github/issues/stlimtat/remiges-smtp.svg?style=for-the-badge
[issues-url]: https://github.com/stlimtat/remiges-smtp/issues
[license-shield]: https://img.shields.io/github/license/stlimtat/remiges-smtp.svg?style=for-the-badge
[license-url]: https://github.com/stlimtat/remiges-smtp/blob/master/LICENSE.txt
