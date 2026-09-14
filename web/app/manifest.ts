import type { MetadataRoute } from "next";

export default function manifest(): MetadataRoute.Manifest {
  return {
    name: "ArkMesh",
    short_name: "ArkMesh",
    description: "Open-source research into verifiable, consent-based preservation of useful AI across infrastructure failure.",
    start_url: "/",
    display: "standalone",
    background_color: "#00072f",
    theme_color: "#001467",
    categories: ["developer", "utilities", "education"],
    icons: [
      {
        src: "/icon.svg",
        sizes: "any",
        type: "image/svg+xml",
      },
    ],
  };
}
