import path from "path";
import type { AliasOptions } from "vite";

export const resolveAlias: AliasOptions = {
    "@components": path.resolve(__dirname, "./src/lib/components"),
    "@lib": path.resolve(__dirname, "./src/lib"),
    "@routes": path.resolve(__dirname, "./src/routes"),
    "@store": path.resolve(__dirname, "./src/lib/store"),
    "@utils": path.resolve(__dirname, "./src/lib/utils"),
    "@views": path.resolve(__dirname, "./src/views"),
    "@layouts": path.resolve(__dirname, "./src/lib/layouts"),
    "@data": path.resolve(__dirname, "./src/data"),
};