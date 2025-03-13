// @ts-check
import { defineConfig } from 'astro/config';

// https://astro.build/config
export default defineConfig({
    server: {
        port: 3030,
    },
    vite: {
        server: {
            watch: {
                usePolling: true, // Force polling mode
                interval: 100, // Adjust interval for better responsiveness
            },
        }
    }

});
