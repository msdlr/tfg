pragma solidity >=0.4.22 <0.7.0;

import "./AuthContract.sol";

contract GeneralContract {

    /* STRUCTS */
    struct User {
            string id;
            bool isRegistered;
            bool isLoggedIn;
            bool adminStatus;
            AuthContract auth;
            uint8 attempts;
    }
    /* EVENTS */
    event createAdmin(address _adm, address _who);

    /* MODIFIERS */
    modifier isUser() {
        require(!userList[msg.sender].isRegistered, "This user is not in the system.");
        _;
    }

    modifier isAdmin() {
        require(!userList[msg.sender].isRegistered, "This user is not in the system.");
        require(userList[msg.sender].adminStatus, "This user does not have admin. priviledges.");
        _;
    }

    modifier userNotLocked{
        require (userList[msg.sender].attempts < 3,"This account is locked");
        _;
    }

    /* Contract data */
    mapping ( address => User)  private userList ;
    mapping ( string => address) private id2a; // Index by id
    address private owner;

    /* CONSTRUCTOR */

    constructor(address owner_, string memory id_) public payable{
        // Set the owner of the company
        owner = owner_;

        // Add it to the admin list
        userList[owner].isRegistered = false;
        userList[owner].adminStatus = true;
        userList[owner].id = id_;
        id2a[id_] = owner;
        
        // Others fields are initialized as default values (0, false)
    }

    /* -- ADMIN FUNCTIONS -- */

    function rmUser(address _addr, string memory _id) public isAdmin {
        require(_addr != owner && id2a[_id] != owner,"You cannot remove the owner");
        // _addr = id2a[_id]
        userList[_addr].auth.terminate();
        userList[_addr].isRegistered = false;
        userList[_addr].adminStatus = false;
        userList[_addr].isLoggedIn = false;
        id2a[_id] = address(0);
        userList[_addr].id = "";
        userList[_addr].attempts = 0;
    }

    function addUser(address _addr, string memory _id) public isAdmin {
        userList[_addr].auth = AuthContract(_addr);
        userList[_addr].isRegistered = false;
        userList[_addr].adminStatus = false;
        userList[_addr].isLoggedIn = false;
        userList[_addr].id = _id;
        id2a[_id] = _addr;
        userList[_addr].attempts = 0;
    }

    function addAdmin(address _addr) public isAdmin {
        // Check that the user is added
        require(userList[_addr].isRegistered == false,"User does not exist.");
        // We update the user's profile with admin status
        userList[_addr].adminStatus = true;
        // We notify in the blockchain who did it
        emit createAdmin(_addr, msg.sender);
    }

    /* -- USER FUNCTIONS (WRAPPERS)-- */
    function getOTP() public isUser userNotLocked returns(uint16 pass_){
        require (userList[msg.sender].isLoggedIn == false, "Only offline users can get for a key");
        // We call that specific contract function
        pass_ = userList[msg.sender].auth.newOTP();
    }

    function tryLogin(uint16 _pass) public isUser userNotLocked {
        // We call that specific contract function
        require(userList[msg.sender].isLoggedIn == false, "You are already logged in");
        
        // Every time an attempt is made, the count is increased
        if(userList[msg.sender].auth.tryLogin(_pass) == true){
            // Successful login
            userList[msg.sender].attempts = 0;
        }
        else {
            // Failed attempt
            userList[msg.sender].attempts++;
        }
    }
    
    function amILocked() public view isUser returns (bool locked){
        locked = (userList[msg.sender].attempts > 3);
    }

    function tryLogout() public isUser {
        // Can only logout if logged in
        require(userList[msg.sender].isLoggedIn == true, "You are not logged in");
        userList[msg.sender].isLoggedIn = false;
    }
    
    /* SETTERS / GETTERS */
    function getOwner() public view returns(address owner_){
        owner_ = owner;
    }
    
    function setOwner(address _newOwner) public isAdmin{
        owner = _newOwner;
    }
    
    //Functions for retrieving the user struct fields
    
    function getUserId(address _addr) public view returns (string memory id_ ){
        id_ = userList[_addr].id;
    }
    
    function getUserRegistered(address _addr) public view returns (bool b){
        b = userList[_addr].isRegistered;
    }
    
    function getUserLoggedIn(address _addr) public view returns (bool b){
        b = userList[_addr].isLoggedIn;
    }
    
    function getUserAdmin(address _addr) public view returns (bool b){
        b = userList[_addr].adminStatus;
    }
    
    function getUserContract(address _addr) public view returns (AuthContract auth_){
        auth_ = userList[_addr].auth;
    }
    
    function getUserAttempts(address _addr) public view returns (uint8 attempts_){
        attempts_ = userList[_addr].attempts;
    }
    
}
