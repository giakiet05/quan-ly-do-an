// Định dạng cơ bản dùng cho các danh sách (Lecturer, Students trong Class)
export interface UserInfo {
    userId: string;
    fullName: string;
    avatar: string; // Chỗ này thường BE trả về URL trực tiếp khi map vào Class
}

// Định dạng chi tiết cho trang Profile cá nhân
export interface UserProfile {
    id: string;
    email: string;
    fullName: string;
    studentCode?: string;
    provider: "local" | "google";
    avatar?: {
        url: string;
        publicId: string;
    };
    createdAt: string;
}

export interface UpdateProfileRequest {
    fullName?: string;
    studentCode?: string;
}

export interface ChangePasswordRequest {
    old_password: string;
    new_password: string;
}