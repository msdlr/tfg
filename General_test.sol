pragma solidity >=0.4.22 <0.7.0;
import "remix_tests.sol"; // this import is automatically injected by Remix.
import "tfg/General.sol";

// File name has to end with '_test.sol', this file can contain more than one testSuite contracts
contract General_test {
    
    GeneralContract testContract;
    //address tester = 0x5B38Da6a701c568545dCfcB03FcB875f56beddC4; // msg.sender
    
    address[] testAddrs = 
    [0xAb8483F64d9C6d1EcF9b849Ae677dD3315835cb2,
    0x4B20993Bc481177ec7E8f571ceCaE8A9e22C02db,
    0x78731D3Ca6b7E34aC0F824c42a7cC18A495cabaB,
    0x617F2E2fD72FD9D5503197092aC168c91465E7f2,
    0x17F6AD8Ef982297579C203069C1DbfFE4348c372,
    0x5c6B0f7Bf3E7ce046039Bd8FABdfD3f9F5021678,
    0x03C6FcED478cBbC9a4FAB34eF9f40767739D1Ff7,
    0x03C6FcED478cBbC9a4FAB34eF9f40767739D1Ff7,
    0x0A098Eda01Ce92ff4A4CCb7A4fFFb5A43EBC70DC,
    0xCA35b7d915458EF540aDe6068dFe2F44E8fa733c,
    0x14723A09ACff6D2A60DcdF7aA4AFf308FDDC160C,
    0x4B0897b0513fdC7C541B6d9D7E929C4e5364D2dB,
    0x583031D1113aD414F02576BD6afaBfb302140225,
    0xdD870fA1b7C4700F2BD7f44238821C26f7392148];
    

    /// 'beforeAll' runs before all other tests
    /// More special functions are: 'beforeEach', 'beforeAll', 'afterEach' & 'afterAll'
    function beforeAll() public {
        // Instantiate the contract to test
        testContract = new GeneralContract(msg.sender, "11223344K");
    }
    
    function checkFirstAdmin() public {
        // Check that the constructor executed correctly
        Assert.equal(msg.sender, testContract.getOwner(), "owner address should be this caller");
        
        // Check if the fields initialized correctly
        Assert.equal(testContract.getUserAdmin(msg.sender),true,"caller was not made admin");
        Assert.equal(testContract.getUserRegistered(msg.sender),false,"caller is not on the user list");
        Assert.equal(testContract.getUserLoggedIn(msg.sender),false,"caller is not on the user list");
        
    }
    
    function addUser() public {
        
        // Try to add a new user
        //Assert.equal(testContract.checkRegistered(testAddrs[0]),false,"User entry (isRegistered) should be FALSE");
        //testContract.addUser(testAddrs[0], "11223344K");
        
        // Check it was done correctly
        //Assert.equal(testContract.checkRegistered(testAddrs[0]),false,"User entry (isRegistered) should be TRUE");
    }
    
    function removeUser() public {
        
    }
    
    
    
}
