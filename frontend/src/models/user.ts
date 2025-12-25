// User model interface for frontend usage
export interface User {
  id: string;
  avatar: string;
  email: string;
  fullname: string;
  role: string;
  is_verified: boolean;
}
