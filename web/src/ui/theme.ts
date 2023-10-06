import { CSSProperties, createTheme, keyframes } from "@macaron-css/core"

const space = {
  "0": "0px",
  px: "1px",
  "0.5": "0.125rem",
  "1": "0.25rem",
  "1.5": "0.375rem",
  "2": "0.5rem",
  "2.5": "0.625rem",
  "3": "0.75rem",
  "3.5": "0.875rem",
  "4": "1rem",
  "5": "1.25rem",
  "6": "1.5rem",
  "7": "1.75rem",
  "8": "2rem",
  "9": "2.25rem",
  "10": "2.5rem",
  "11": "2.75rem",
  "12": "3rem",
  "14": "3.5rem",
  "16": "4rem",
  "20": "5rem",
  "24": "6rem",
  "28": "7rem",
  "32": "8rem",
  "36": "9rem",
  "40": "10rem",
  "44": "11rem",
  "48": "12rem",
  "52": "13rem",
  "56": "14rem",
  "60": "15rem",
  "64": "16rem",
  "72": "18rem",
  "80": "20rem",
  "96": "24rem",
};

const base = {
  radius: "0.5rem",
}

const lightBase = {
  background: "0, 0%, 100%",
  foreground: "222.2, 84%, 4.9%",

  nav: "0, 0%, 98%",
  navForeground: "222.2, 84%, 6.9%",

  card: "0, 0%, 100%",
  cardForeground: "222.2, 84%, 4.9%",

  popover: "0, 0%, 100%",
  popoverForeground: "222.2, 84%, 4.9%",

  primary: "221.2, 83.2%, 53.3%",
  primaryForeground: "210, 40%, 98%",

  secondary: "210, 40%, 96.1%",
  secondaryForeground: "222.2, 47.4%, 11.2%",

  muted: "210, 40%, 96.1%",
  mutedForeground: "215.4, 16.3%, 46.9%",

  accent: "210, 40%, 96.1%",
  accentForeground: "222.2, 47.4%, 11.2%",

  destructive: "0, 84.2%, 60.2%",
  destructiveForeground: "210, 40%, 98%",

  border: "214.3, 31.8%, 91.4%",
  input: "214.3, 31.8%, 91.4%",
  ring: "221.2, 83.2%, 53.3%",
};

const light = {
  ...base,
  background: `hsl(${lightBase.background})`,
  foreground: `hsl(${lightBase.foreground})`,
  backgroundOverlay: `hsla(${lightBase.background}, 80%)`,

  nav: `hsl(${lightBase.nav})`,
  navForeground: `hsl(${lightBase.navForeground})`,

  card: `hsl(${lightBase.card})`,
  cardForeground: `hsl(${lightBase.cardForeground})`,

  popover: `hsl(${lightBase.popover})`,
  popoverForeground: `hsl(${lightBase.popoverForeground})`,

  primary: `hsl(${lightBase.primary})`,
  primaryForeground: `hsl(${lightBase.primaryForeground})`,
  primaryHover: `hsla(${lightBase.primary}, 90%)`,
  primaryOverlay: `hsla(${lightBase.primary}, 80%)`,

  secondary: `hsl(${lightBase.secondary})`,
  secondaryForeground: `hsl(${lightBase.secondaryForeground})`,
  secondaryHover: `hsla(${lightBase.secondary}, 80%)`,

  muted: `hsl(${lightBase.muted})`,
  mutedForeground: `hsl(${lightBase.mutedForeground})`,
  mutedHover: `hsla(${lightBase.muted}, 50%)`,

  accent: `hsl(${lightBase.accent})`,
  accentForeground: `hsl(${lightBase.accentForeground})`,

  destructive: `hsl(${lightBase.destructive})`,
  destructiveForeground: `hsl(${lightBase.destructiveForeground})`,
  destructiveHover: `hsla(${lightBase.destructive}, 90%)`,
  destructiveOverlay: `hsla(${lightBase.destructive}, 80%)`,

  border: `hsl(${lightBase.border})`,
  input: `hsl(${lightBase.input})`,
  ring: `hsl(${lightBase.ring})`,
};

const darkBase = {
  background: "222.2, 84%, 4.9%",
  foreground: "210, 40%, 98%",

  nav: "222.2, 84%, 6.9%",
  navForeground: "210, 40%, 96%",

  card: "222.2, 84%, 4.9%",
  cardForeground: "210, 40%, 98%",

  popover: "222.2, 84%, 4.9%",
  popoverForeground: "210, 40%, 98%",

  primary: "217.2, 91.2%, 59.8%",
  primaryForeground: "222.2, 47.4%, 11.2%",

  secondary: "217.2, 32.6%, 17.5%",
  secondaryForeground: "210, 40%, 98%",

  muted: "217.2, 32.6%, 17.5%",
  mutedForeground: "215, 20.2%, 65.1%",

  accent: "217.2, 32.6%, 17.5%",
  accentForeground: "210, 40%, 98%",

  destructive: "0, 62.8%, 30.6%",
  destructiveForeground: "210, 40%, 98%",

  border: "217.2, 32.6%, 17.5%",
  input: "217.2, 32.6%, 17.5%",
  ring: "224.3, 76.3%, 48%",
};

const dark = {
  ...base,
  background: `hsl(${darkBase.background})`,
  foreground: `hsl(${darkBase.foreground})`,
  backgroundOverlay: `hsla(${darkBase.background}, 80%)`,

  nav: `hsl(${darkBase.nav})`,
  navForeground: `hsl(${darkBase.navForeground})`,

  card: `hsl(${darkBase.card})`,
  cardForeground: `hsl(${darkBase.cardForeground})`,

  popover: `hsl(${darkBase.popover})`,
  popoverForeground: `hsl(${darkBase.popoverForeground})`,

  primary: `hsl(${darkBase.primary})`,
  primaryForeground: `hsl(${darkBase.primaryForeground})`,
  primaryHover: `hsla(${darkBase.primary}, 90%)`,
  primaryOverlay: `hsla(${darkBase.primary}, 80%)`,

  secondary: `hsl(${darkBase.secondary})`,
  secondaryForeground: `hsl(${darkBase.secondaryForeground})`,
  secondaryHover: `hsla(${darkBase.secondary}, 80%)`,

  muted: `hsl(${darkBase.muted})`,
  mutedForeground: `hsl(${darkBase.mutedForeground})`,
  mutedHover: `hsla(${darkBase.muted}, 50%)`,

  accent: `hsl(${darkBase.accent})`,
  accentForeground: `hsl(${darkBase.accentForeground})`,

  destructive: `hsl(${darkBase.destructive})`,
  destructiveForeground: `hsl(${darkBase.destructiveForeground})`,
  destructiveHover: `hsla(${darkBase.destructive}, 90%)`,
  destructiveOverlay: `hsla(${darkBase.destructive}, 80%)`,

  border: `hsl(${darkBase.border})`,
  input: `hsl(${darkBase.input})`,
  ring: `hsl(${darkBase.ring})`,
};

const size = {
  sm: "640px",
  md: "768px",
  lg: "1024px",
  xl: "1280px",
  "2xl": "1536px",
};

const borderRadius = {
  // None: "0px",
  sm: "0.125rem",
  ok: "0.25rem",
  md: "0.375rem",
  lg: "0.5rem",
  // xl: "0.75rem",
  // 2xl: "1rem",
  // 3xl: "1.5rem",
  full: "9999px"
}

export const minScreen = {
  sm: "screen and (min-width: 640px)",
  md: "screen and (min-width: 768px)",
  lg: "screen and (min-width: 1024px)",
  xl: "screen and (min-width: 1280px)",
  "2xl": "screen and (min-width: 1536px)",
};

const themeDefault = {
  space,
  size,
  borderRadius,
};

export const [darkClass, theme] = createTheme({
  ...themeDefault,
  color: {
    ...dark,
  },
});

export const lightClass = createTheme(theme, {
  ...themeDefault,
  color: {
    ...light,
  },
});


const rotate = keyframes({
  from: { transform: "rotate(0deg)" },
  to: { transform: "rotate(360deg)" },
});

export const tw = {
  animateSpin: {
    animation: `${rotate} 1s linear infinite`,
  } as CSSProperties,

  shadowSm: {
    boxShadow: "0 1px 2px 0 rgb(0 0 0 / 0.05)"
  } as CSSProperties,

  shadow: {
    boxShadow: "0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1)"
  } as CSSProperties,

  shadowMd: {
    boxShadow: "0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1);"
  } as CSSProperties,

  shadowLg: {
    boxShadow: "0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1);"
  } as CSSProperties,

  // shadowXl: {
  //   boxShadow: "0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1);"
  // } as CSSProperties,

  // shadow2x: {
  //   boxShadow: "0 25px 50px -12px rgb(0 0 0 / 0.25);"
  // } as CSSProperties,

  // shadowInner: {
  //   boxShadow: "inset 0 2px 4px 0 rgb(0 0 0 / 0.05);"
  // } as CSSProperties,

  // shadowNone: {
  //   boxShadow: "0 0 #0000"
  // } as CSSProperties,

  textXs: {
    fontSize: "0.75rem",
    lineHeight: "1rem"
  } as CSSProperties,

  textSm: {
    fontSize: "0.875rem",
    lineHeight: "1.25rem"
  } as CSSProperties,

  textOk: {
    fontSize: "1rem",
    lineHeight: "1.5rem"
  } as CSSProperties,

  textLg: {
    fontSize: "1.125rem",
    lineHeight: "1.75rem"
  } as CSSProperties,

  textXl: {
    fontSize: "1.25rem",
    lineHeight: "1.75rem"
  } as CSSProperties,

  text2xl: {
    fontSize: "1.5rem",
    lineHeight: "2rem"
  } as CSSProperties,

  transitionColors: {
    transitionProperty: "color, background-color, border-color, text-decoration-color, fill, stroke",
    transitionTimingFunction: "cubic-bezier(0.4, 0, 0.2, 1)",
    transitionDuration: "150ms"
  } as CSSProperties
}

export const mixin = {
  textLine(): CSSProperties {
    return {
      overflow: "hidden",
      textOverflow: "ellipsis",
      whiteSpace: "nowrap",
    };
  },

  stack(space: keyof typeof theme["space"]): CSSProperties {
    return {
      display: "flex",
      flexDirection: "column",
      gap: theme.space[space],
    };
  },

  row(space: keyof typeof theme["space"]): CSSProperties {
    return {
      display: "flex",
      gap: theme.space[space],
    };
  },

  size(space: keyof typeof theme["space"]): CSSProperties {
    return {
      width: theme.space[space],
      height: theme.space[space],
    };
  },
};
