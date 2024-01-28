const clientKey = '%s'
const serverURL = '%s'

document.addEventListener('DOMContentLoaded', function () {
  var utmData = getUTMParameters()

  if (utmData) {
    console.log('UTM data found')
    sendUTMData(utmData)
  }

  // Send page view data to the server
  var pageURL = window.location.href
  sendPageViewData(pageURL)
})

// Send page view data to the server navigation change
window.addEventListener('hashchange', function () {
  var pageURL = window.location.href
  sendPageViewData(pageURL)
})

document.addEventListener('click', function (event) {
  var allowedTags = ['button', 'a', 'span']

  if (!allowedTags.includes(event.target.tagName.toLowerCase())) {
    return
  }
  var clickedElement = {
    tag: event.target.tagName.toLowerCase(),
    id: event.target.id,
    classList: Array.from(event.target.classList),
    textContent: event.target.textContent.trim(),
  }

  if (clickedElement.tag === 'a') {
    clickedElement.href = event.target.href
  }

  var parentElement = {
    tag: event.target.parentElement.tagName.toLowerCase(),
    id: event.target.parentElement.id,
    classList: Array.from(event.target.parentElement.classList),
    textContent: event.target.textContent.trim(),
  }

  if (parentElement.tag === 'a') {
    parentElement.href = event.target.parentElement.href
  }

  clickedElement.parentElement = parentElement

  var pageURL = window.location.href

  body = {
    element: clickedElement,
    url: pageURL,
  }

  sendClickData(body)
})

// Function to parse UTM parameters from the URL
function getUTMParameters() {
  var queryParams = new URLSearchParams(window.location.search)
  var utmSource = queryParams.get('utm_source')
  var utmCampaign = queryParams.get('utm_campaign')
  var utmMedium = queryParams.get('utm_medium')
  var track = queryParams.get('track')

  if (!utmSource && !utmCampaign && !utmMedium && !track) {
    return null
  }

  return {
    utm_source: utmSource,
    utm_campaign: utmCampaign,
    utm_medium: utmMedium,
    track: track,
  }
}

// Function to send UTM data to the tracking server
function sendUTMData(utmData) {
  var pageURL = window.location.href.split('?')[0]

  body = {
    ...utmData,
    page_url: pageURL,
  }

  fetch(serverURL + '/api/v1/track/utm', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Site-Key': clientKey,
      Origin: window.location.origin,
    },
    body: JSON.stringify({
      ...body,
    }),
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error('Failed to send UTM data to the server')
      }
    })
    .catch((error) => {
      console.error(error)
    })
}

// Function to send page view data to the tracking server
function sendPageViewData(pageURL) {
  fetch(serverURL + '/api/v1/track/pageview', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Site-Key': clientKey,
      Origin: window.location.origin,
    },
    body: JSON.stringify({
      url: pageURL,
    }),
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error('Failed to send page view data to the server')
      }
    })
    .catch((error) => {
      console.error(error)
    })
}

// Function to send click data to the tracking server
function sendClickData(data) {
  fetch(serverURL + '/api/v1/track/click', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Site-Key': clientKey,
      Origin: window.location.origin,
    },
    body: JSON.stringify(data),
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error('Failed to send click data to the server')
      }
    })
    .catch((error) => {
      console.error(error)
    })
}
