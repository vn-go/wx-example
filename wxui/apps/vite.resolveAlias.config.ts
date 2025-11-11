// src/vite/base.config.ts
import path from 'path';

export const resolveAlias = {
    '@components': path.resolve('./src/lib/components'),
    '@lib': path.resolve('./src/lib'),
    '@routes': path.resolve('./src/routes'),
    '@store': path.resolve('./src/lib/store'),
    '@utils': path.resolve('./src/lib/utils'),
    '@layouts': path.resolve('./src/layouts'),
};