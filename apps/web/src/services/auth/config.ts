
const isHttps = window.location.protocol === "https:";

export const COOKIE_OPTIONS: Cookies.CookieAttributes = {
    expires: 7,
    secure: isHttps,
    sameSite: "lax",
    path: "/"
};