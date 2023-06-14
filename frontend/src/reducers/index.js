import user from "./user";
import { combineReducers } from "redux";
import { connect as reduxConnect } from "react-redux";
import { clearAuth } from "./../reducers/auth"
import { updateUser } from "./../reducers/user"

const rootReducer = combineReducers({
  user
});

function mapStateToProps(state) {
  const user = state.user.value
  const auth = state.auth.value
  return {
    auth,
    user
  }
}

function mapDispatchToProps(dispatch) {
  return {
    clearAuth: () => dispatch(clearAuth()),
    updateUser: user => dispatch(updateUser(user))
  }
}

// connect can be used to replace the react-redux connect call, and provides default maps for state and dispatch
export const connect = (component) => {
  return reduxConnect(mapStateToProps, mapDispatchToProps)(component)
}

export default rootReducer;