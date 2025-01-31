import { defineStore } from 'pinia';
import type { IUser } from 'src/models/user/user';
import { AuthService } from 'src/services/auth';
import { IResource } from 'src/services/resource/resource';
import { TokenService } from 'src/utils/tokenService';
import { Ref, ref } from 'vue';
import { useRouter } from 'vue-router';
import { jwtDecode } from "jwt-decode"; // Make sure to install this library: `npm install jwt-decode`

export const useAuthStore = defineStore('auth', () => {
  const state: Ref<IResource | null> = ref(null);
  const userType: Ref<IUser.UserType | null> = ref(
    TokenService.type as IUser.UserType
  );
  const router = useRouter();

  const login = async (data: IUser.Login) => {
    const res = await AuthService.login(data);
    state.value = res;
    userType.value = res.data ?? null;
  };

  const tryApplyJwt = (token: string): boolean => {
    try {
      const decoded: { exp: number; iat: number; user_id: number; user_post: string } = jwtDecode(token);
      const currentTime = Math.floor(Date.now() / 1000); // Current time in seconds
      if (decoded.exp < currentTime) {
        console.error("Token is expired", currentTime, decoded.exp);
        return false;
      }
      userType.value = decoded.user_post as IUser.UserType;
      TokenService.token = token;
      return true;
    } catch (error) {
      console.error("Failed to apply JWT:", error);
      return false;
    }
  };

  const register = async (data: IUser.Register) => {
    const res = await AuthService.register(data);
    state.value = res;
  };

  const logout = async () => {
    TokenService.token = '';
    TokenService.type = '';
    userType.value = null;
    state.value = null;
    router.push('/');
  };

  const forgetPassword = async (email: string) => {
    await AuthService.forgetPassword({ email });
  };

  const restorePassword = async (password: string, token: string) => {
    const res = await AuthService.restorePassword({
      new_password: password,
      token,
    });
    if (!res.error) {
      router.push({ path: '/auth' });
    }
  };

  return {
    state,
    userType,
    login,
    register,
    logout,
    forgetPassword,
    restorePassword,
    tryApplyJwt
  };
});
