// Error codes cho Classroom operations (theo API_DOCUMENTATION.md)

export const CLASSROOM_ERROR_CODES = {
  // Join Classroom errors
  CODE_INVALID: "CODE_INVALID",
  ALREADY_IN_CLASSROOM: "ALREADY_IN_CLASSROOM",
  CLASSROOM_FULL: "CLASSROOM_FULL",
  STUDENT_CODE_NOT_IN_WHITELIST: "STUDENT_CODE_NOT_IN_WHITELIST",
  INVALID_EMAIL_DOMAIN: "INVALID_EMAIL_DOMAIN",
  STUDENT_CODE_ALREADY_USED: "STUDENT_CODE_ALREADY_USED",
  
  // Classroom management errors
  CLASSROOM_NOT_FOUND: "CLASSROOM_NOT_FOUND",
  JOIN_REQUEST_NOT_FOUND: "JOIN_REQUEST_NOT_FOUND",
  JOIN_REQUEST_ALREADY_PROCESSED: "JOIN_REQUEST_ALREADY_PROCESSED",
  
  // Permission errors
  FORBIDDEN: "FORBIDDEN",
  UNAUTHORIZED: "UNAUTHORIZED",
} as const;

export type ClassroomErrorCode = typeof CLASSROOM_ERROR_CODES[keyof typeof CLASSROOM_ERROR_CODES];

// Error messages hiển thị cho user (tiếng Việt)
export const CLASSROOM_ERROR_MESSAGES: Record<ClassroomErrorCode, string> = {
  [CLASSROOM_ERROR_CODES.CODE_INVALID]: 
    "Mã mời không hợp lệ hoặc đã hết hạn",
  [CLASSROOM_ERROR_CODES.ALREADY_IN_CLASSROOM]: 
    "Bạn đã là thành viên của lớp học này",
  [CLASSROOM_ERROR_CODES.CLASSROOM_FULL]: 
    "Lớp học đã đầy, không thể tham gia",
  [CLASSROOM_ERROR_CODES.STUDENT_CODE_NOT_IN_WHITELIST]: 
    "Mã sinh viên của bạn không có trong danh sách được phép tham gia",
  [CLASSROOM_ERROR_CODES.INVALID_EMAIL_DOMAIN]: 
    "Email của bạn không thuộc domain được phép tham gia lớp này",
  [CLASSROOM_ERROR_CODES.STUDENT_CODE_ALREADY_USED]: 
    "Mã sinh viên này đã được sử dụng bởi người dùng khác",
  [CLASSROOM_ERROR_CODES.CLASSROOM_NOT_FOUND]: 
    "Không tìm thấy lớp học",
  [CLASSROOM_ERROR_CODES.JOIN_REQUEST_NOT_FOUND]: 
    "Không tìm thấy yêu cầu tham gia",
  [CLASSROOM_ERROR_CODES.JOIN_REQUEST_ALREADY_PROCESSED]: 
    "Yêu cầu này đã được xử lý",
  [CLASSROOM_ERROR_CODES.FORBIDDEN]: 
    "Bạn không có quyền thực hiện thao tác này",
  [CLASSROOM_ERROR_CODES.UNAUTHORIZED]: 
    "Vui lòng đăng nhập để tiếp tục",
};

// Helper function để lấy error message
export function getClassroomErrorMessage(errorCode: string): string {
  return CLASSROOM_ERROR_MESSAGES[errorCode as ClassroomErrorCode] 
    || "Đã xảy ra lỗi không xác định";
}
