// vue.config.js
const webpack = require('webpack');
const isProd = process.env.NODE_ENV === "production";

module.exports = {
    outputDir: "../backend/static",
    configureWebpack: {
        devServer: {
            proxy: {
                "/ws": {
                    "target": "ws://localhost:8000",
                    "ws": true,
                    "secure": false,
                    "changeOrigin": true
                },
            }
        }
    }
};

