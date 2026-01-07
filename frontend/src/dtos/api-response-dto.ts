/**
 * Standard API response wrapper
 */
export interface ApiResponse<T = any> {
  success: boolean;
  message: string;
  data?: T;
  error_code?: string;
}

/**
 * Paginated response wrapper
 */
export interface PaginatedResponse<T = any> {
  data: T[];
  pagination: {
    current_page: number;
    total_pages: number;
    total_items: number;
    page_size: number;
  };
}

/**
 * Error response
 */
export interface ErrorResponse {
  success: false;
  message: string;
  error_code: string;
}
