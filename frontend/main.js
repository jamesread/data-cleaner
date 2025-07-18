import { createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';

import { DataCleanerService } from './resources/javascript/gen/data_cleaner/api/v1/data_cleaner_pb';

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

  const tbl = dom.querySelector('tbody')

  for (const sourceFile of summary.sourceFiles) {
    const row = document.createElement('tr')

    const filename = document.createElement('td')
    filename.innerText = sourceFile.filename
    row.appendChild(filename)

    const lineCount = document.createElement('td')
    lineCount.innerText = sourceFile.lineCount
    row.appendChild(lineCount)

    tbl.appendChild(row)
  }
}

function createApiClient () {
  const transport = createConnectTransport({
    baseUrl: './api',
  })

  window.apiClient = createClient(DataCleanerService, transport)
}

export async function main () {
  createApiClient()

  try { 
    const res = await window.apiClient.import()

    if (res.issues.length === 0) {
      const msg = document.createElement('p')
      msg.classList.add('inline-notification')
      msg.classList.add('good')
      msg.innerText = 'No issues found!'

      document.getElementById('issuesList').appendChild(msg)

      document.getElementById('loadJobButton').removeAttribute('disabled')
      document.getElementById('loadJobButton').onclick = () => {
        loadJob()
      }
    } else {
      displayIssues(res.issues)
    }

    displaySourceFiles(res)
    displaySummary(res)
    displayTransformations(res.transformations)
  } catch (error) {
    showBigError('Failed to import: ' + error.message)
  }
}

function displayTransformations (transformations) {
  const dom = document.getElementById('transformStatus')

  for (const transformation of transformations) {
    const li = document.createElement('li')

    li.innerText = transformation.description
    dom.appendChild(li)
  }
}

function loadJob () {
  const loadButton = document.getElementById('loadJobButton')
  loadButton.setAttribute('disabled', true)
  loadButton.classList.remove('good')
  loadButton.innerText = 'Loading';

  console.log(loadButton)

  try {
    const response = window.apiClient.load()

    loadButton.classList.add('good')
    loadButton.innerText = 'Loaded successfully!';

    setTimeout(() => {
      loadButton.removeAttribute('disabled')
      loadButton.innerText = 'Load';
    }, 1000)

    console.log('Job loaded successfully:', response)
  } catch (error) {
    showBigError('Failed to load: ' + error.message)
  }
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
