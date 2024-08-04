Deployment: https://rms-be-482r.onrender.com

## Routes
1.  `POST /signup`
    - Create a profile on the system (Name, Email, Password, UserType (Admin/Applicant), Profile Headline, Address).
2. `POST /login`
    - Authenticate users and return a JWT token upon successful validation.
3. `POST /uploadResume`
    - Authenticated API for uploading resume files (only PDF or DOCX) of the applicant. Only Applicant type users can access this API.
4. `POST /admin/job``
    - Authenticated API for creating job openings. Only Admin type users can access this API.
5. `GET /admin/job/{job_id}`
    - Authenticated API for fetching information regarding a job opening.
    - Returns details about the job opening and a list of applicants. Only Admin type users can access this API.
6. `GET /admin/applicants`
    - Authenticated API for fetching a list of all users in the system. Only Admin type users can access this API.
7. `GET /admin/applicant/{applicant_id}`
    - Authenticated API for fetching extracted data of an applicant. Only Admin type users can access this API.
8. `GET /jobs`
    - Authenticated API for fetching job openings. All users can access this API.
9. `GET /jobs/apply?job_id={job_id}`
    - Authenticated API for applying to a particular job. Only Applicant users are allowed to apply for jobs.
    
10. `GET /score/{job_id}`
    -  Authenticated API for getting match scores for thier particular job.
    -  Only Applicant users are allowed get match scores
