import { createSlice } from "@reduxjs/toolkit";

const initialState = {
    value: {
        status: "",
        signOutStatus: "",
    }
}

const authSlice = createSlice({
    name: "auth",
    initialState,
    reducers: {
        setAuthPending(state) {
            state.value.status = "PENDING"
        },
        setAuthComplete(state) {
            state.value.status = "COMPLETE"
        },
        clearAuth(state) {
            state.value.status = ""
        },
        setSignOutPending(state) {
            state.value.signOutStatus = "PENDING"
        },
        setSignOutComplete(state) {
            state.value.signOutStatus = "COMPLETE"
        },
        clearSignOut(state) {
            state.value.signOutStatus = ""
        }
    }
})

export const { setAuthPending, setAuthComplete, clearAuth, setSignOutPending, setSignOutComplete, clearSignOut } = authSlice.actions
export default authSlice.reducer