export interface CreateClassPostRequest {
    title: string;
    content: string;
    attachments?: Array<AttachmentUpload>;
}
// {{base_url}}/api/classrooms/{{classroom_id}}/posts POST
// {
//     "title": "Thông báo nghỉ học",
//     "content": "Lớp nghỉ ngày 10/1/2026",
//     "attachments": [
//         {
//             "file_name": "schedule.pdf",
//             "file_url": "https://cloudinary.com/...",
//             "file_size": 1024000,
//             "mime_type": "application/pdf"
//         }
//     ]
// }

