# Build stage
FROM node:18-alpine AS builder

# Accept build arguments for Next.js public environment variables
ARG NEXT_PUBLIC_API_BASE_URL
ARG NEXT_PUBLIC_NEW_RELIC_BROWSER_KEY
ARG NEXT_PUBLIC_NEW_RELIC_ACCOUNT_ID
ARG NEXT_PUBLIC_NEW_RELIC_APPLICATION_ID

# Set them as environment variables for the build
ENV NEXT_PUBLIC_API_BASE_URL=$NEXT_PUBLIC_API_BASE_URL
ENV NEXT_PUBLIC_NEW_RELIC_BROWSER_KEY=$NEXT_PUBLIC_NEW_RELIC_BROWSER_KEY
ENV NEXT_PUBLIC_NEW_RELIC_ACCOUNT_ID=$NEXT_PUBLIC_NEW_RELIC_ACCOUNT_ID
ENV NEXT_PUBLIC_NEW_RELIC_APPLICATION_ID=$NEXT_PUBLIC_NEW_RELIC_APPLICATION_ID

WORKDIR /app

# Install dependencies needed for native modules
RUN apk add --no-cache libc6-compat

# Copy package files
COPY package*.json ./
RUN npm ci

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage
FROM node:18-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# Create a non-root user
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Copy necessary files from builder
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules

# Change ownership
RUN chown -R nextjs:nodejs /app

USER nextjs

EXPOSE 3000

ENV PORT=3000

CMD ["npm", "start"]