# URL Shortener
##  Create compact and user-friendly short links from lengthy URLs.

![Diagram](/assets/urlshortener.png)
Developed a URL shortener project with a primary focus on generating concise and short URLs. It includes implementation of Redis caching to optimize performance. Moreover, a rate-limiting feature has been integrated to enhance both security and scalability aspects of the system.

### Built using
1. Go
2. Redis
3. NGINX
4. MongoDB
5. Docker
6. GitHub Actions
7. AWS ECR

### Key Features:
1. **Expiry for Short URLs:** Easily create short URLs with optional expiration dates. After the specified expiry, the short link becomes invalid. If no expiry is provided during link generation, a default expiry of 1 year is applied.

2. **Efficient Caching and Database Integration:** Store recently used short URLs in cache for rapid access. If a URL is absent from cache, the system seamlessly checks the MongoDB database. Cached URLs are returned and stored, while missing URLs trigger database queries, with results both stored in cache and returned to users. 

3. **Smart Cache Management:** Employ the `maxmemory-policy allkeys-lfu` configuration for Redis caching. When memory reaches its limit, the system intelligently removes the least used URLs, ensuring space for frequently accessed ones and optimizing memory usage.

4. **Spam Protection with Nginx Rate Limiting:** Bolster system integrity by implementing rate limiting through an Nginx reverse proxy. This safeguard prevents spam requests from generating random URLs, effectively reducing server load and maintaining overall performance.

5. **Data Consistency**: By periodically running the Database Purge API, We can guarantees that only valid and up-to-date short URLs remain accessible. This API is not accessible by user.

6. **AWS ECR**: Utilized GitHub Actions to automate the building and pushing of Docker images to AWS Elastic Container Registry (ECR) whenever code is pushed to the main branch, ensuring swift and accurate deployment with version-tagged images.

### Project Setup
Running the project is straightforward. Follow these steps:
1. Ensure you have Docker installed on your system.
2. Open a terminal window and navigate to the project directory.
3. Run the following command to start the project using Docker Compose:
   ```
   docker compose up
   ```
   This command will set up and launch the required services, making your project accessible.

4. That's it! Your project is now up and running.
   Open your web browser and access the project using the provided URLs or endpoints.

Enjoy using the project!

### Demo
Generate URL
- url : Add webpage url
- expiry: Add expiry time (in Hour)

<img src="/assets/short1.png" width="640px"/>

AWS ECR

<img src="/assets//ecr.png" width="640px"/>




