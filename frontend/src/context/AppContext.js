import { createContext, useContext, useEffect, useReducer } from "react";
import {
  clearToken as clearClientToken,
  createItem as apiCreateItem,
  deleteItem as apiDeleteItem,
  fetchItems as apiFetchItems,
  login as apiLogin,
  register as apiRegister,
  setToken as setClientToken,
  updateItem as apiUpdateItem,
} from "../api/client";

const AppContext = createContext(undefined);

const storage = typeof window !== "undefined" ? window.localStorage : null;

const storedToken = storage?.getItem("app_token") || null;
let storedUser = null;
if (storage) {
  try {
    storedUser = JSON.parse(storage.getItem("app_user") || "null");
  } catch {
    storedUser = null;
  }
}

const initialState = {
  user: storedUser,
  token: storedToken,
  items: [],
  loading: false,
  error: null,
  notification: null,
};

function reducer(state, action) {
  switch (action.type) {
    case "SET_LOADING":
      return { ...state, loading: action.payload };
    case "SET_ERROR":
      return { ...state, error: action.payload };
    case "SET_NOTIFICATION":
      return { ...state, notification: action.payload };
    case "LOGIN_SUCCESS":
      return {
        ...state,
        user: action.payload.user,
        token: action.payload.token,
        error: null,
      };
    case "LOGOUT":
      return { ...state, user: null, token: null, items: [] };
    case "SET_ITEMS":
      return { ...state, items: action.payload };
    case "ADD_ITEM":
      return { ...state, items: [action.payload, ...state.items] };
    case "UPDATE_ITEM":
      return {
        ...state,
        items: state.items.map((item) =>
          item.id === action.payload.id ? action.payload : item
        ),
      };
    case "REMOVE_ITEM":
      return {
        ...state,
        items: state.items.filter((item) => item.id !== action.payload),
      };
    default:
      return state;
  }
}

export function AppProvider({ children }) {
  const [state, dispatch] = useReducer(reducer, initialState);

  useEffect(() => {
    if (state.token) {
      setClientToken(state.token);
      storage?.setItem("app_token", state.token);
      storage?.setItem("app_user", JSON.stringify(state.user));
    } else {
      clearClientToken();
      storage?.removeItem("app_token");
      storage?.removeItem("app_user");
    }
  }, [state.token, state.user]);

  useEffect(() => {
    if (state.token) {
      fetchItems();
    } else {
      dispatch({ type: "SET_ITEMS", payload: [] });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [state.token]);

  function setLoading(isLoading) {
    dispatch({ type: "SET_LOADING", payload: isLoading });
  }

  function setError(message) {
    dispatch({ type: "SET_ERROR", payload: message });
  }

  function setNotification(message) {
    dispatch({ type: "SET_NOTIFICATION", payload: message });
  }

  async function login(credentials) {
    setLoading(true);
    setError(null);
    try {
      const data = await apiLogin(credentials);
      dispatch({ type: "LOGIN_SUCCESS", payload: data });
      setNotification("Signed in successfully");
      return true;
    } catch (error) {
      const message =
        error.response?.data?.error || "Unable to sign in, please try again.";
      setError(message);
      return false;
    } finally {
      setLoading(false);
    }
  }

  async function register(credentials) {
    setLoading(true);
    setError(null);
    try {
      await apiRegister(credentials);
      setNotification("Registration completed. Please sign in.");
      return true;
    } catch (error) {
      const message =
        error.response?.data?.error ||
        "Unable to complete registration. Please try again.";
      setError(message);
      return false;
    } finally {
      setLoading(false);
    }
  }

  function logout() {
    dispatch({ type: "LOGOUT" });
    setNotification("You have been signed out.");
  }

  async function fetchItems() {
    if (!state.token) {
      return false;
    }
    setLoading(true);
    setError(null);
    try {
      const items = await apiFetchItems();
      dispatch({ type: "SET_ITEMS", payload: items });
      return true;
    } catch (error) {
      const message =
        error.response?.data?.error || "Unable to load items right now.";
      setError(message);
      return false;
    } finally {
      setLoading(false);
    }
  }

  async function createItem(payload) {
    setLoading(true);
    setError(null);
    try {
      const item = await apiCreateItem(payload);
      dispatch({ type: "ADD_ITEM", payload: item });
      setNotification("Item created successfully.");
      return true;
    } catch (error) {
      const message =
        error.response?.data?.error || "Unable to create item.";
      setError(message);
      return false;
    } finally {
      setLoading(false);
    }
  }

  async function updateItem(id, payload) {
    setLoading(true);
    setError(null);
    try {
      const item = await apiUpdateItem(id, payload);
      dispatch({ type: "UPDATE_ITEM", payload: item });
      setNotification("Item updated successfully.");
      return true;
    } catch (error) {
      const message =
        error.response?.data?.error || "Unable to update item.";
      setError(message);
      return false;
    } finally {
      setLoading(false);
    }
  }

  async function deleteItem(id) {
    setLoading(true);
    setError(null);
    try {
      await apiDeleteItem(id);
      dispatch({ type: "REMOVE_ITEM", payload: id });
      setNotification("Item removed.");
      return true;
    } catch (error) {
      const message =
        error.response?.data?.error || "Unable to remove item.";
      setError(message);
      return false;
    } finally {
      setLoading(false);
    }
  }

  const value = {
    state,
    actions: {
      login,
      register,
      logout,
      fetchItems,
      createItem,
      updateItem,
      deleteItem,
      setError,
      setNotification,
      setLoading,
    },
  };

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>;
}

export function useAppContext() {
  const context = useContext(AppContext);
  if (!context) {
    throw new Error("useAppContext must be used within AppProvider");
  }
  return context;
}

