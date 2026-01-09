// src/dtos/auth-dto.ts
import type { User } from '../models/user';


export interface LoginRequestDto {
  identifier: string; // username hoặc email
  password: string;
}

export interface LoginResponseDto {
  user: User;
  access_token: string;
  refresh_token: string;
}


export interface SendVerificationRequestDto {
  email: string;
}


export interface VerifyEmailRequestDto {
  email: string;
  otp: string;
}

export interface VerifyEmailResponseDto {
  verification_token: string;
}


export interface CompleteRegistrationRequestDto {
  verification_token: string;
  full_name: string;
  password: string;
}

export interface RegisterResponseDto {
  user: User;
  access_token: string;
  refresh_token: string;
}


export interface RefreshTokenRequestDto {
  refresh_token: string;
}

export interface RefreshTokenResponseDto {
  access_token: string;
  refresh_token: string;
}


export interface LogoutRequestDto {
  access_token: string;
  refresh_token: string;
}



export interface ForgotPasswordRequestDto {
  email: string;
}

export interface VerifyResetOtpRequestDto {
  email: string;
  otp: string;
}

export interface VerifyResetOtpResponseDto {
  reset_token: string;
}

export interface ResetPasswordRequestDto {
  reset_token: string;
  new_password: string;
}


export interface CompleteGoogleSetupRequestDto {
  setup_token: string;
  username: string;
}
