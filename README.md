# Simple Site Tracker

## Overview

The Simple Site Tracker is a lightweight and easy-to-use tool designed to be a self hosted site traffic tracker based on UTM parameters, page clicks and page views.
It offers a simple setup process with a single endpoint for generating new site integrations, and the app provides an endpoint which serves the minified JS.

## Why ?
I wanted to have analytics on my small web apps, but I didn't need the complexity of Google Analytics, and free versions of alternatives were very restrictive. 

This app provides me an API to easily build my own dashboard, and I have direct access to the data. So I can run SQL queries to get out any specific information I need.

Self hosting this tracker also allows me to notify myself on certain events. (ie. Tracking when people from a certain utm_campaign viewed my site.)

## Features

- **UTM Parameter Tracking:** Keep tabs on UTM campaigns associated with your URLs to understand the effectiveness of various marketing efforts.

- **Page Clicks:** Track user clicks on different pages of your website to gain insights into user engagement.

- **Page Views:** Monitor the overall page views to assess the popularity and performance of your website.

- **JavaScript Generation:** Easy integration with a simple JavaScript snippet. Users only need to add the provided script to their web pages.

- **Validation:** Validation included to ensure that only your domain can be tracked against, which helps against malicious actors.

## Usage

1. Include the following script in the `<head>` section of your HTML file:

   ```html
   <script src="https://appurl/server/js/{clientKey}"></script>
   ```




