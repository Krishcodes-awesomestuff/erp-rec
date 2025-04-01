function generateCaptcha() {
    const digits = '0123456789';
    let captcha = '';
    for(let i = 0; i < 4; i++) {  // Changed to exactly 4 digits
        captcha += digits[Math.floor(Math.random() * 10)] + ' ';
    }
    return captcha.trim();
}