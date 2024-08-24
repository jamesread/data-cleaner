const moneyFormatter = new Intl.NumberFormat('en-GB', {
  style: 'currency',
  currency: 'GBP',
})

function newIssue () {
  return document.getElementById('issueTemplate').content.cloneNode(true)
}

function displayIssues (issues) {
  const list = document.getElementById('issuesList')

  for (const issue of issues) {
    const dom = newIssue()

    dom.querySelector('.description').innerText = issue.description

    const detailsTable = dom.querySelector('.detailsTable')

    renderRowAttributes(detailsTable, issue.expected, 'Expected')
    renderRowAttributes(detailsTable, issue.intermediate, 'Intermediate')
    renderRowAttributes(detailsTable, issue.actual, 'Actual')

    dom.querySelector('.filename').innerText = issue.locationFilename
    dom.querySelector('.lineNumber').innerText = issue.locationLineNumber

    list.appendChild(dom)
  }
}

function renderRowAttributes (detailsTable, rowAttributes, title) {
  const row = document.createElement('tr')

  const expectedLbl = document.createElement('td')
  expectedLbl.innerText = title + ':'
  row.appendChild(expectedLbl)

  for (const expected of rowAttributes) {
    const expectedKey = document.createElement('td')
    expectedKey.innerText = expected.key
    row.appendChild(expectedKey)

    const expectedVal = document.createElement('td')
    expectedVal.innerText = expected.val
    row.appendChild(expectedVal)

    detailsTable.appendChild(row)
  }
}

function main () {
  window.fetch('./api/Import')
    .then(response => {
      if (!response.ok) {
        showBigError('Failed to fetch issues: ' + response.statusText)
      } else {
        return response.json()
      }
    })
    .then(res => {
      if (res.issues.length === 0) {
        const msg = document.createElement('li')
        msg.innerText = 'No issues found!'

        document.getElementById('issuesList').appendChild(msg)
      } else {
        displayIssues(res.issues)
      }
    })
    .catch(error => showBigError(error))
}

function showBigError (error) {
  const msg = document.createElement('dialog')
  msg.classList.add('alert')
  msg.classList.add('critical')
  msg.innerText = 'An error occurred: ' + error

  document.body.appendChild(msg)

  msg.showModal()

  console.log(error)

  throw new Error(error)
}

main()
