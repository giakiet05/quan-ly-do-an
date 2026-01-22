export const AppErrorCode = {
    // Auth-related
    INVALID_CREDENTIALS: "INVALID_CREDENTIALS",
    INVALID_TOKEN: "INVALID_TOKEN",
    INVALID_CLAIMS: "INVALID_CLAIMS",
    INVALID_ISSUER: "INVALID_ISSUER",
    INVALID_AUDIENCE: "INVALID_AUDIENCE",
    TOKEN_INVALIDATED: "TOKEN_INVALIDATED",
    FORBIDDEN: "FORBIDDEN",
    BAD_REQUEST: "BAD_REQUEST",
    EMAIL_NOT_VERIFIED: "EMAIL_NOT_VERIFIED",
    EMAIL_ALREADY_VERIFIED: "EMAIL_ALREADY_VERIFIED",
    INVALID_OTP: "INVALID_OTP",
    OTP_EXPIRED: "OTP_EXPIRED",
    LOGIN_METHOD_MISMATCH: "LOGIN_METHOD_MISMATCH",
    EMAIL_NOT_REGISTERED: "EMAIL_NOT_REGISTERED",

    // Generic
    INTERNAL_ERROR: "INTERNAL_ERROR",
    NO_FIELDS_TO_UPDATE: "NO_FIELDS_TO_UPDATE",
    INVALID_ID: "INVALID_ID",
    PAGINATION_INVALID: "PAGINATION_INVALID",
    NOT_FOUND: "NOT_FOUND",

    // User-related
    USER_NOT_FOUND: "USER_NOT_FOUND",
    USERNAME_EXISTS: "USERNAME_EXISTS",
    EMAIL_EXISTS: "EMAIL_EXISTS",
    USER_INACTIVE: "USER_INACTIVE",

    // Project-related
    PROJECT_NOT_FOUND: "PROJECT_NOT_FOUND",
    REPORT_ALREADY_EXISTS: "REPORT_ALREADY_EXISTS",
    PROJECT_GROUP_LIMIT_REACHED: "PROJECT_GROUP_LIMIT_REACHED",

    // Classroom-related
    CLASSROOM_NOT_FOUND: "CLASSROOM_NOT_FOUND",
    ALREADY_IN_CLASSROOM: "ALREADY_IN_CLASSROOM",
    CLASSROOM_FULL: "CLASSROOM_FULL",
    CODE_INVALID: "CODE_INVALID",
    STUDENT_CODE_NOT_IN_WHITELIST: "STUDENT_CODE_NOT_IN_WHITELIST",
    STUDENT_CODE_ALREADY_USED: "STUDENT_CODE_ALREADY_USED",
    INVALID_EMAIL_DOMAIN: "INVALID_EMAIL_DOMAIN",

    // Invitation-related
    INVITATION_NOT_FOUND: "INVITATION_NOT_FOUND",
    INVITATION_EXPIRED: "INVITATION_EXPIRED",
    INVITATION_ALREADY_PROCESSED: "INVITATION_ALREADY_PROCESSED",
} as const;

/**
 * Map các mã lỗi sang thông báo Tiếng Việt
 */
export const AppErrorMessage: Record<string, string> = {
    // Auth
    [AppErrorCode.INVALID_CREDENTIALS]: "Tên đăng nhập hoặc mật khẩu không chính xác.",
    [AppErrorCode.INVALID_TOKEN]: "Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.",
    [AppErrorCode.FORBIDDEN]: "Bạn không có quyền thực hiện hành động này.",
    [AppErrorCode.EMAIL_NOT_VERIFIED]: "Email của bạn chưa được xác thực.",
    [AppErrorCode.LOGIN_METHOD_MISMATCH]: "Email này đã được đăng ký bằng phương thức khác. Vui lòng sử dụng đúng cách đăng nhập ban đầu.",
    [AppErrorCode.EMAIL_NOT_REGISTERED]: "Email này chưa được đăng ký trong hệ thống.",
    [AppErrorCode.INVALID_OTP]: "Mã xác thực (OTP) không đúng.",
    [AppErrorCode.OTP_EXPIRED]: "Mã xác thực đã hết hạn.",

    // Generic & User
    [AppErrorCode.INTERNAL_ERROR]: "Lỗi hệ thống. Vui lòng thử lại sau.",
    [AppErrorCode.NOT_FOUND]: "Không tìm thấy dữ liệu yêu cầu.",
    [AppErrorCode.USER_NOT_FOUND]: "Người dùng không tồn tại.",
    [AppErrorCode.EMAIL_EXISTS]: "Email này đã được sử dụng.",
    [AppErrorCode.USER_INACTIVE]: "Tài khoản của bạn hiện đang bị khóa.",

    // Classroom & Student
    [AppErrorCode.CLASSROOM_NOT_FOUND]: "Không tìm thấy lớp học.",
    [AppErrorCode.ALREADY_IN_CLASSROOM]: "Bạn đã là thành viên của lớp học này rồi.",
    [AppErrorCode.CLASSROOM_FULL]: "Lớp học đã đủ số lượng thành viên.",
    [AppErrorCode.CODE_INVALID]: "Mã tham gia lớp học không hợp lệ.",
    [AppErrorCode.STUDENT_CODE_NOT_IN_WHITELIST]: "Mã số sinh viên của bạn không nằm trong danh sách cho phép của lớp này.",
    [AppErrorCode.STUDENT_CODE_ALREADY_USED]: "Mã số sinh viên này đã được sử dụng để đăng ký lớp học này.",
    [AppErrorCode.INVALID_EMAIL_DOMAIN]: "Vui lòng sử dụng email sinh viên (@gm.uit.edu.vn) để tham gia.",

    // Project
    [AppErrorCode.PROJECT_NOT_FOUND]: "Không tìm thấy dự án.",
    [AppErrorCode.PROJECT_GROUP_LIMIT_REACHED]: "Nhóm dự án đã đủ thành viên.",
};

/**
 * Hàm lấy thông báo lỗi từ mã code (Dùng cho cả URL params và API response)
 */
export function getErrorMessage(code: string | null | undefined): string {
    if (!code) return "Đã xảy ra lỗi không xác định.";

    // Ưu tiên lấy trong map tiếng Việt, nếu không thấy thì trả về chính mã code đó để debug
    return AppErrorMessage[code] || `Lỗi hệ thống: ${code}`;
}