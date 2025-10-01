# Imgix Clone - Development Roadmap

## Epic 1: Project Foundation & Setup
**Goal:** Get basic project structure and serve your first image

### Task 1.1: Project Initialization
- [x] Initialize Go module (`go mod init`)
- [x] Choose and set up project structure (Standard Go project layout)
- [x] Set up `.gitignore` and basic documentation
- [x] Choose HTTP framework (or use standard library)

### Task 1.2: Basic HTTP Server
- [x] Create main.go with basic HTTP server
- [x] Add health check endpoint (`GET /health`)
- [x] Set up graceful shutdown
- [x] Test server runs and responds

### Task 1.3: Configuration Management
- [ ] Create config structure (port, upload directory, etc.)
- [ ] Load config from environment variables
- [ ] Add validation for config values

**Deliverable:** Running HTTP server with health check endpoint

---

## Epic 2: Image Upload & Storage
**Goal:** Accept and store images locally

### Task 2.1: File Upload Endpoint
- [ ] Create `POST /upload` endpoint
- [ ] Parse multipart/form-data
- [ ] Validate file is an image (check MIME type)
- [ ] Generate unique filename (UUID or hash)

### Task 2.2: Local Storage Implementation
- [ ] Create storage interface for future flexibility
- [ ] Implement local filesystem storage
- [ ] Save uploaded images to disk
- [ ] Return image URL/ID in response

### Task 2.3: Image Retrieval
- [ ] Create `GET /images/:id` endpoint
- [ ] Serve original images from disk
- [ ] Handle proper Content-Type headers
- [ ] Add error handling for missing files

**Deliverable:** Can upload and retrieve images

---

## Epic 3: Basic Image Transformations
**Goal:** Resize images on-the-fly

### Task 3.1: Image Processing Library Setup
- [ ] Research libraries (imaging, bimg, disintegration/imaging)
- [ ] Install chosen library
- [ ] Create image processing service/package
- [ ] Test basic image load and save

### Task 3.2: Resize Transformation
- [ ] Parse query parameters (`?w=300&h=200`)
- [ ] Implement resize logic
- [ ] Handle aspect ratio preservation
- [ ] Return transformed image

### Task 3.3: Query Parameter Parsing
- [ ] Create parameter parser middleware/helper
- [ ] Validate parameter values
- [ ] Set defaults for missing params
- [ ] Handle invalid inputs gracefully

**Deliverable:** Can resize images via URL parameters

---

## Epic 4: Additional Transformations
**Goal:** Support multiple image operations

### Task 4.1: Format Conversion
- [ ] Add `format` parameter (jpg, png, webp)
- [ ] Implement format conversion
- [ ] Update Content-Type based on format
- [ ] Handle format-specific options

### Task 4.2: Quality Control
- [ ] Add `quality` parameter (1-100)
- [ ] Implement quality adjustment
- [ ] Set sensible defaults per format

### Task 4.3: Crop Operations
- [ ] Add `fit` parameter (crop, fill, scale)
- [ ] Implement different crop modes
- [ ] Add `crop` parameter for manual cropping
- [ ] Handle crop positioning

**Deliverable:** Multiple transformation options working

---

## Epic 5: Caching Layer
**Goal:** Improve performance with caching

### Task 5.1: Cache Strategy Design
- [ ] Design cache key structure (URL + params)
- [ ] Choose cache backend (in-memory first, Redis later)
- [ ] Define cache expiration policy
- [ ] Plan cache invalidation strategy

### Task 5.2: In-Memory Cache Implementation
- [ ] Implement simple in-memory cache (map + sync.RWMutex)
- [ ] Add cache middleware to image serving
- [ ] Implement LRU eviction (use library or simple approach)
- [ ] Add cache hit/miss logging

### Task 5.3: Cache Headers
- [ ] Add proper Cache-Control headers
- [ ] Implement ETag generation
- [ ] Handle conditional requests (304 Not Modified)

**Deliverable:** Cached image serving with better performance

---

## Epic 6: URL Signing & Security
**Goal:** Secure image URLs and prevent abuse

### Task 6.1: URL Signature Generation
- [ ] Design signature format (HMAC-SHA256)
- [ ] Create signing function
- [ ] Add signature to URLs
- [ ] Create CLI tool or endpoint for generating signed URLs

### Task 6.2: Signature Verification
- [ ] Add verification middleware
- [ ] Validate signatures on requests
- [ ] Return 403 for invalid signatures
- [ ] Add timestamp expiration (optional)

### Task 6.3: Rate Limiting
- [ ] Implement basic rate limiting per IP
- [ ] Add rate limit headers
- [ ] Handle rate limit exceeded responses

**Deliverable:** Signed URLs preventing unauthorized access

---

## Epic 7: Cloud Storage Integration
**Goal:** Store images in S3 instead of local disk

### Task 7.1: Storage Interface Refactoring
- [ ] Create storage interface (if not done)
- [ ] Refactor existing code to use interface
- [ ] Ensure local storage still works

### Task 7.2: S3 Client Setup
- [ ] Install AWS SDK for Go
- [ ] Configure S3 credentials
- [ ] Create S3 storage implementation
- [ ] Test upload and download

### Task 7.3: Migration & Configuration
- [ ] Add config for storage backend selection
- [ ] Support both local and S3 storage
- [ ] Update upload endpoint for S3
- [ ] Update retrieval endpoint for S3

**Deliverable:** Images stored in S3 with same API

---

## Epic 8: Advanced Features
**Goal:** Add professional-grade features

### Task 8.1: Image Optimization
- [ ] Auto-optimize images on upload
- [ ] Strip metadata (EXIF) option
- [ ] Progressive JPEG support
- [ ] WebP auto-conversion for supported browsers

### Task 8.2: Advanced Transformations
- [ ] Blur effect
- [ ] Grayscale filter
- [ ] Brightness/contrast adjustments
- [ ] Rotation

### Task 8.3: Monitoring & Metrics
- [ ] Add Prometheus metrics
- [ ] Track transformation times
- [ ] Monitor cache hit rates
- [ ] Log slow operations

**Deliverable:** Production-ready image service

---

## Epic 9: Testing & Documentation
**Goal:** Ensure quality and usability

### Task 9.1: Unit Tests
- [ ] Test image processing functions
- [ ] Test parameter parsing
- [ ] Test cache logic
- [ ] Aim for >70% coverage

### Task 9.2: Integration Tests
- [ ] Test full upload flow
- [ ] Test transformation pipeline
- [ ] Test error scenarios

### Task 9.3: Documentation
- [ ] API documentation (endpoints, parameters)
- [ ] Setup instructions
- [ ] Architecture diagrams
- [ ] Example usage

**Deliverable:** Well-tested and documented project

---

## Epic 10: Deployment
**Goal:** Run in production

### Task 10.1: Dockerization
- [ ] Create Dockerfile
- [ ] Multi-stage build for smaller image
- [ ] Docker compose for local development

### Task 10.2: Deployment Preparation
- [ ] Environment-based configuration
- [ ] Logging setup (structured logging)
- [ ] Error reporting
- [ ] Health checks

### Task 10.3: Deploy
- [ ] Choose platform (Fly.io, Railway, AWS)
- [ ] Deploy application
- [ ] Set up domain
- [ ] Test production environment

**Deliverable:** Live image processing service!

---

## Recommended Order

Start with: **Epic 1 → Epic 2 → Epic 3**

This gives you a working MVP where you can upload and resize images.

Then tackle epics based on interest:
- Want performance? → Epic 5 (Caching)
- Want security? → Epic 6 (URL Signing)
- Want cloud? → Epic 7 (S3)
- Want polish? → Epic 8 (Advanced Features)