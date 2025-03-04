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
smtp client with file scrapping



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
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

The project was conceived as a smtpclient with limited capabilities:
1. Read a set of files from a directory
2. Parse the files and format them as RFC compliant emails
  a. Sign the message with DKIM-Signature
3. Connect to the mail servers indicated by the To address and send the email
4. Write the result in a CSV file

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

To get a local copy up and running follow these simple example steps.

### Prerequisites

To get started, you will need to have downloaded the dependent libraries.  This has been managed by Bazel.

* smtpclient cli tool
  ```sh
  bazel test //...
  ```
* unit tests
  ```sh
  bazel test //...
  ```

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/stlimtat/remiges-smtp.git
   ```
2. Run the simple smtpclient server in a docker container
   ```sh
   docker compose up
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

See [USAGE.md](./doc/USAGE.md)

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

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Top contributors:

<a href="https://github.com/stlimtat/remiges-smtp/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=stlimtat/remiges-smtp" alt="contrib.rocks image" />
</a>

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
