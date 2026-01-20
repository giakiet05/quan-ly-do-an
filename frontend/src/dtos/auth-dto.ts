// src/dtos/auth-dto.ts
import type { User } from '../models/user';


export interface LoginRequestDto {
  identifier: string;
  password: string;
}

export interface LoginResponseDto {
  user: User;
  accessToken: string;
  refreshToken: string;
}


export interface SendVerificationRequestDto {
  email: string;
}


export interface VerifyEmailRequestDto {
  email: string;
  otp: string;
}

export interface VerifyEmailResponseDto {
  verificationToken: string;
}


export interface CompleteRegistrationRequestDto {
  verificationToken: string;
  fullName: string;
  password: string;
}

export interface RegisterResponseDto {
  user: User;
  accessToken: string;
  refreshToken: string;
}


export interface RefreshTokenRequestDto {
  refreshToken: string;
}

export interface RefreshTokenResponseDto {
  accessToken: string;
  refreshToken: string;
}


export interface LogoutRequestDto {
  accessToken: string;
  refreshToken: string;
}



export interface ForgotPasswordRequestDto {
  email: string;
}

export interface VerifyResetOtpRequestDto {
  email: string;
  otp: string;
}

export interface VerifyResetOtpResponseDto {
  resetToken: string;
}

export interface ResetPasswordRequestDto {
  resetToken: string;
  newPassword: string;
}


export interface CompleteGoogleSetupRequestDto {
  setupToken: string;
  fullName: string;
}
