pragma solidity >=0.4.22 <0.7.0;

contract AUTH {
// <editor-fold defaultstate="collapsed" desc="Structs">
    struct User {
        address uAddr; // Used for checking a null user
        bool uLoggedIn;
        bool isAdmin;
    }

    struct OTP {
        User user;
        uint16 pass;
        uint256 timestamp;
        bool isUsed;
        bool isExpired;
    }
// </editor-fold>

// <editor-fold defaultstate="collapsed" desc="State variables">
    mapping ( address => User) userList;
    mapping ( address => OTP) otpList;

    // The approximate time for the OTPs to expire
    uint constant  OTPtimeout = 5 minutes;

// </editor-fold>

// <editor-fold defaultstate="collapsed" desc="Events">
    event createAdmin(address _adm, address _who, uint256 t);
// </editor-fold>

// <editor-fold defaultstate="collapsed" desc="Modifiers">

    modifier isUser() {
        require(userList[msg.sender].uAddr != address(0), "This user is not in the system.");
        _;
    }

    modifier isAdmin() {
        require(userList[msg.sender].uAddr != address(0), "This user is not in the system.");
        require(userList[msg.sender].isAdmin, "This user does not have admin. priviledges.");
        _;
    }

    modifier unusedOTP(){
        // The OTP can only be used once
        require(otpList[msg.sender].isUsed == false, "OTP is already used");
        _;
    }
// </editor-fold>

// <editor-fold defaultstate="collapsed" desc="Constructor">

    constructor(address _creator) public payable{
        // Called by the master account
        _creator = msg.sender;

        // Stablish the parameter address as the first admin
        userList[_creator].uAddr = msg.sender;
        userList[_creator].isAdmin = true;
    }
// </editor-fold>

// <editor-fold defaultstate="collapsed" desc="Admin functions">
    function rmUser(address _addr) public isAdmin {
        userList[_addr].uAddr = address(0);
        userList[_addr].isAdmin = false;
    }

    function addUser(address _addr) public isAdmin {
        userList[_addr].uAddr = address(0);
        userList[_addr].isAdmin = false;
    }

    function addAdmin(address _addr) public isAdmin {
        // Check that the user is added
        require(userList[_addr].uAddr != address (0),"User does not exist.");
        // We update the user's profile with admin status
        userList[_addr].isAdmin = true;
        // We notify in the blockchain who did it
        emit createAdmin(_addr, msg.sender, now);
    }
// </editor-fold>

// <editor-fold defaultstate="collapsed" desc="User functions">
    function userLoggedIn(address _addr) public view isUser returns (bool b){
        if(userList[_addr].isAdmin){
            // Admins can check for every user
            return (userList[_addr].uLoggedIn);
        }
        else {
            // Users cannot check for other users
            require(userList[_addr].uAddr == msg.sender,"You can only check yourself.");
            return userList[msg.sender].uLoggedIn;
        }
    }
// </editor-fold>

// <editor-fold defaultstate="collapsed" desc="OTP functions">

    function genOTP() private {
        // We generate a random number from 0 to 9999
        uint16 pass = uint16(uint256( keccak256( abi.encode(now, msg.sender) ) ) % 9999);
    
        // We fill the fields for the OTP
        otpList[msg.sender].user = userList[msg.sender];
        otpList[msg.sender].timestamp = block.timestamp;
        otpList[msg.sender].pass = pass;
        otpList[msg.sender].isUsed = false;
        otpList[msg.sender].isExpired = false;
    }

    // This is the pass and time remaining
    // It returns the caller's OTP
    function getOTP() public isUser view returns
    (uint16 OTPnumber, uint aproxTime){
        // We only need to return the OTP pass
        // And an aproximate remaining time
        // Returns negative value if expired
        aproxTime = (now + OTPtimeout - otpList[msg.sender].timestamp);
        return (otpList[msg.sender].pass, aproxTime);
    }

    function checkOTPexpired(address _addr) private view returns(bool b) {
        // Returns if the OTP has expired
        return (now < (OTPtimeout + otpList[_addr].timestamp));
    }
// </editor-fold>

}
