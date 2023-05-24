function randomEmail() {
    return emails[Math.floor(Math.random() * emails.length)]
}

const emails = [
    "example@gmail.com"
]

export default randomEmail