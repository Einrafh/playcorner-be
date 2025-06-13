import type { ZudokuConfig } from "zudoku";

const config: ZudokuConfig = {
  page: {
    logo: {
      src: { light: "/logo-light.svg", dark: "/logo-dark.svg" },
      alt: "PlayCorner API",
      width: "150px",
    },
  },
  navigation: [
    {
      type: "category",
      label: "PlayCorner Docs",
      items: [
        {
          type: "category",
          label: "Panduan",
          icon: "sparkles",
          items: [
            "/introduction",
          ],
        },
        {
          type: "link",
          label: "API Reference",
          icon: "folder-cog",
          to: "/api",
        },
      ],
    },
    {
      type: "link",
      to: "/api",
      label: "Full API Reference",
    },
  ],
  redirects: [{ from: "/", to: "/introduction" }],
  apis: [
    {
      type: "file",
      input: "./apis/playcorner-api.yaml",
      path: "/api",
    },
  ],
};

export default config;

