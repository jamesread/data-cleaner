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

    const hdr = document.createElement('tr')
    createAppend(hdr, 'th', '')

    for (const attr of issue.expected) {
      createAppend(hdr, 'th', attr.key)
    }
    detailsTable.appendChild(hdr)

    renderRowAttributes(detailsTable, issue.expected, 'Expected')
    renderRowAttributes(detailsTable, issue.actual, 'Actual')
    renderRowAttributes(detailsTable, issue.intermediate, 'Intermediate')

    dom.querySelector('.filename').innerText = issue.currentLocationFilename
    dom.querySelector('.lineNumber').innerText = issue.currentLocationLineNumber
    dom.querySelector('.lastFilename').innerText = issue.lastLocationFilename
    dom.querySelector('.lastLineNumber').innerText = issue.lastLocationLineNumber


    list.appendChild(dom)
  }
}

function createAppend (parent, tag, text) {
  const elem = document.createElement(tag)
  elem.innerText = text

  parent.appendChild(elem)
}

function renderRowAttributes (detailsTable, rowAttributes, title) {
  const row = document.createElement('tr')

  const expectedLbl = document.createElement('td')
  expectedLbl.classList.add('uneditable')
  expectedLbl.innerText = title + ':'
  row.appendChild(expectedLbl)

  for (const expected of rowAttributes) {
    const expectedVal = document.createElement('td')
    expectedVal.innerText = expected.val
    row.appendChild(expectedVal)
  }

  detailsTable.appendChild(row)
}

function displaySummary (summary) {
  document.getElementById('totalLines').innerText = summary.totalLines
  document.getElementById('totalFiles').innerText = summary.sourceFiles.length
}

function displaySourceFiles (summary) {
  const dom = document.getElementById('sourceFiles')

  for (const sourceFile of summary.sourceFiles) {
    const file = document.createElement('li')
    file.innerText = sourceFile.filename + ' (' + sourceFile.lineCount + ' lines)'
    dom.appendChild(file)
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
        displaySourceFiles(res)
        displaySummary(res)
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
