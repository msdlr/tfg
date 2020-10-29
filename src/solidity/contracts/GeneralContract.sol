pragma solidity >=0.6.4 <=7.3.0;

import "./AuthContract.sol";
import "./GenericSensorContract.sol";

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
    event promoteEvent(address _adm, address _who);
    event demoteEvent(address _adm, address _who);

    /* MODIFIERS */
    modifier onlyRegistered() {
        require(userList[msg.sender].isRegistered, "Caller is not in the system.");
        _;
    }

    modifier onlyAdmin() {
        require(userList[msg.sender].isRegistered, "Caller is not in the system.");
        require(userList[msg.sender].adminStatus, "Caller does not have admin. priviledges.");
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
        userList[owner].isRegistered = true;
        userList[owner].adminStatus = true;
        userList[owner].id = id_;
        id2a[id_] = owner;
        userList[owner].auth = new AuthContract(this, owner);
        // Others fields are initialized as default values (0, false)
    }

    /* -- ADMIN FUNCTIONS -- */

    function rmUser(address _addr, string memory _id) public onlyAdmin {
        require(_addr != owner && id2a[_id] != owner,"You cannot remove the owner");
        require(id2a[_id] == _addr,"Address and Id do not correspond");
        // _addr = id2a[_id]
        userList[_addr].auth.terminate();
        userList[_addr].isRegistered = false;
        userList[_addr].adminStatus = false;
        userList[_addr].isLoggedIn = false;
        id2a[_id] = address(0);
        userList[_addr].id = "";
        userList[_addr].attempts = 0;
    }

    function addUser(address _addr, string memory _id) public onlyAdmin {
        require(userList[_addr].isRegistered == false,"User is already registered.");
        require(id2a[_id] == address(0),"Id is already taken");
        userList[_addr].auth = new AuthContract(this, _addr);
        userList[_addr].isRegistered = true;
        userList[_addr].adminStatus = false;
        userList[_addr].isLoggedIn = false;
        userList[_addr].id = _id;
        id2a[_id] = _addr;
        userList[_addr].attempts = 0;
    }

    function promoteUser(address _addr) public onlyAdmin {
        // Check that the user is added
        require(userList[_addr].isRegistered == true,"User does not exist.");
        // We update the user's profile with admin status
        userList[_addr].adminStatus = true;
        // We notify in the blockchain who did it
        emit promoteEvent(_addr, msg.sender);
    }

    function demoteAdmin(address _addr) public onlyAdmin {
        // Check that the user is added
        require(userList[_addr].isRegistered == true,"User does not exist.");
        // We update the user's profile with admin status
        userList[_addr].adminStatus = false;
        // We notify in the blockchain who did it
        emit demoteEvent(_addr, msg.sender);
    }

    /* -- USER FUNCTIONS (WRAPPERS)-- */
    function getOTP() public onlyRegistered userNotLocked returns(uint16 pass_){
        require (userList[msg.sender].isLoggedIn == false, "Only offline users can get for a key");
        // We call that specific contract function
        pass_ = userList[msg.sender].auth.newOTP();
    }

    function tryLogin(uint16 _pass) public onlyRegistered userNotLocked {
        // We call that specific contract function
        require(userList[msg.sender].isLoggedIn == false, "You are already logged in");
        // Every time an attempt is made, the count is increased
        if(userList[msg.sender].auth.tryLogin(_pass) == true){
            // Successful login
            userList[msg.sender].attempts = 0;
            userList[msg.sender].isLoggedIn = true;
        }
        else {
            // Failed attempt
            userList[msg.sender].attempts++;
        }
    }

    function amILocked() public view onlyRegistered returns (bool locked){
        locked = (userList[msg.sender].attempts > 3);
    }

    function tryLogout() public onlyRegistered {
        // Can only logout if logged in
        require(userList[msg.sender].isLoggedIn == true, "You are not logged in");
        userList[msg.sender].isLoggedIn = false;
    }

    /* SETTERS / GETTERS */
    function getOwner() public view returns(address owner_){
        return owner;
    }

    function setOwner(address _newOwner) public{
        require(msg.sender == owner,"Only the owner can do this");
        owner = _newOwner;
    }

    //Functions for retrieving the user struct fields
    function getContractAddress() public view returns (address){
        return address(this);
    }

    function getUserId(address _addr) public view returns (string memory){
        return userList[_addr].id;
    }

    function getUserAddress(string memory _id) public view returns (address){
        return id2a[_id];
    }

    function getUserRegistered(address _addr) public view returns (bool){
        return userList[_addr].isRegistered;
    }

    function getUserLoggedIn(address _addr) public view returns (bool){
        return userList[_addr].isLoggedIn;
    }

    function getUserAdmin(address _addr) public view returns (bool){
        return userList[_addr].adminStatus;
    }

    function getUserAuthContract(address _addr) public view returns (AuthContract){
        return userList[_addr].auth;
    }

    function getUserAttempts(address _addr) public view returns (uint8){
        return userList[_addr].attempts;
    }
}
