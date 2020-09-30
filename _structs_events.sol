pragma solidity >=0.4.22 <0.7.0;

struct User {
        bool isNull;
        bool isLoggedIn;
        bool isAdmin;
    }

    struct OTP {
        User user;
        uint16 pass;
        uint256 timestamp;
        bool isUsed;
        bool isExpired;
}